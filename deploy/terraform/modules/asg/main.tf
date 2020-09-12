variable "vpc" {}
variable "public_subnets" {}
variable "internet_gateway" {}
variable "account_id" {}

# The behavior when "faas_ami" changes is:
# (1) New LC is created with the fresh AMI
# (2) New ASG is created with the fresh LC
# (3) Terraform waits for the new ASG's instances to spin up and attach to ELB
# (4) Once all new instances are InService, Terraform begins destroy of old ASG
# (5) Once old ASG is destroyed, Terraform destroys old LC
# If Terraform hits its 10m timeout during (3), the new ASG will be marked as
# "tainted" and the apply will halt, leaving the old ASG in service.
# https://groups.google.com/g/terraform-tool/c/7Gdhv1OAc80/m/iNQ93riiLwAJ?pli=1

data "aws_ami" "faas_ami" {
	owners		= [var.account_id]
	most_recent	= true

	filter {
		name	= "name"
		values	= ["faas-http-*"]
	}
}

# lc's cannot be updated, unique naming assures most recent version used
resource "aws_launch_configuration" "faas_server" {
	# omit the "name" attribute to allow Terraform to auto-generate random
	name_prefix					= "faas-httpd-server-"
	image_id					= data.aws_ami.faas_ami.id
	iam_instance_profile		= aws_iam_instance_profile.faas_service.name
	instance_type				= "t2.micro"
	associate_public_ip_address	= true

	security_groups				= [
		aws_security_group.ssh_in.id,
		aws_security_group.http_in.id,
		aws_security_group.http_out.id,
		aws_security_group.ephemeral_out.id,
		aws_security_group.https_out.id,
	]

	# rolling deployments: create & verify new lc before removing old
	lifecycle {
		create_before_destroy	= true
	}
}

resource "aws_autoscaling_group" "faas_service" {
	# interpolate launch configuration name into its name
	# so LC changes always force replacement of the ASG, not just update
	name					= "${aws_launch_configuration.faas_server.name}-asg"
	# Terraform will wait for instances in the new ASG to show up as
	# InService in the ELB before considering the ASG successfully created. 

	min_elb_capacity		= 1
	min_size				= 1
	desired_capacity		= 2
	max_size				= 4
	launch_configuration	= aws_launch_configuration.faas_server.name
	health_check_type		= "ELB"
	vpc_zone_identifier		= var.public_subnets[*].id

	# rolling deployments: create & verify new asg before removing old
	lifecycle {
		create_before_destroy	= true
	}

	enabled_metrics			= [
		"GroupPendingInstances",
		"GroupStandbyInstances",
		"GroupInServiceInstances",
		"GroupTerminatingInstances",
		"GroupTotalInstances",
	]

	# Terraform currently provides both a standalone ASG Attachment resource
	# (describing an ASG attached to an ELB or ALB), and an AutoScaling Group
	# resource with load_balancers and target_group_arns defined in-line.
	# At this time you can use an ASG with in-line load balancers or
	# target_group_arns in conjunction with an ASG Attachment resource,
	# however, to prevent unintended resource updates, the
	# aws_autoscaling_group resource must be configured to ignore changes to the
	# load_balancers and target_group_arns arguments within a lifecycle
	# configuration block.
	target_group_arns		= [
		aws_lb_target_group.faas_target_group.arn,
	]

	tags					= [
		{
			key					= "faas"
			value				= "service"
			propagate_at_launch	= true
		},
		{
			key					= "faas"
			value				= "autoscaling-group"
			propagate_at_launch	= false
		},
	]
}

# resource "aws_s3_bucket" "lb_access_logs_s3_bucket" {
#   bucket = "lb-access-logs-s3-bucket-${var.account_id}"

#   server_side_encryption_configuration {
#     rule {
#       apply_server_side_encryption_by_default {
#         sse_algorithm = "AES256"
#       }
#     }
#   }
# }

resource "aws_lb" "faas_lb" {
	name								= "lb-faas" 
	internal							= false
	subnets								= var.public_subnets[*].id
	enable_cross_zone_load_balancing	= true

	security_groups	= [
		aws_security_group.http_in.id,
		aws_security_group.http_out.id,
		aws_security_group.ephemeral_out.id,
		aws_security_group.https_out.id,
	]

	# Access logs are created only if the load balancer has a TLS
	# listener and they contain information only about TLS requests.
	# https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-access-logs.html
	# access_logs {
	# 	bucket	= aws_s3_bucket.lb_access_logs_s3_bucket.bucket
	# 	enabled	= false
	# }

	tags	=	{
		faas = "lb"
	}
}

output "lb_dns_name" {
	value	= aws_lb.faas_lb.dns_name
}

# assign specific ports to receive incoming traffic
resource "aws_lb_listener" "lb_http_listener" {
	load_balancer_arn	= aws_lb.faas_lb.arn
	port				= "80"
	protocol			= "HTTP"

	default_action {
		target_group_arn	= aws_lb_target_group.faas_target_group.arn
		type				= "forward"
	}
}

resource "aws_lb_listener_rule" "lb_forward_all_rule" {
	listener_arn	= aws_lb_listener.lb_http_listener.arn
	
	action {
		type				= "forward"
		target_group_arn	= aws_lb_target_group.faas_target_group.arn
	}

	condition {
		path_pattern {
			values	= ["*"]
		}
	}
}

# endpoint of LB architecture, receiving forwarded traffic that matched rules
resource "aws_lb_target_group" "faas_target_group" {
	name		= "faas-alb-target-group"
	port		= 80
	protocol	= "HTTP"
	vpc_id		= var.vpc.id
	depends_on	= [aws_lb.faas_lb]

	# Because the health checks act independently if you are using an ASG
	# inside a Target Group configuring them differently can make it difficult
	# to track down where an issue lies.
	# https://medium.com/cognitoiq/terraform-and-aws-application-load-balancers-62a6f8592bcf
	health_check {    
		path		= "/"
		port		= 80
		protocol	= "HTTP"
		matcher		= "200"
	}

	tags		= {
		faas	= "lb_target_group"
	}
}

resource "aws_autoscaling_policy" "autopolicy_down" {
	name					= "faas-autopolicy-down"
	scaling_adjustment		= -1
	adjustment_type			= "ChangeInCapacity"
	cooldown				= 300
	autoscaling_group_name	= aws_autoscaling_group.faas_service.name
}
