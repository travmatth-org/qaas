variable "vpc" {}
variable "public_subnets" {}
variable "internet_gateway" {}
variable "account_id" {}

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
