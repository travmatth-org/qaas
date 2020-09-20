# The behavior when "qaas_ami" changes is:
# (1) New LC is created with the fresh AMI
# (2) New ASG is created with the fresh LC
# (3) Terraform waits for the new ASG's instances to spin up and attach to ELB
# (4) Once all new instances are InService, Terraform begins destroy of old ASG
# (5) Once old ASG is destroyed, Terraform destroys old LC
# If Terraform hits its 10m timeout during (3), the new ASG will be marked as
# "tainted" and the apply will halt, leaving the old ASG in service.
# https://groups.google.com/g/terraform-tool/c/7Gdhv1OAc80/m/iNQ93riiLwAJ?pli=1

data "aws_ami" "qaas_ami" {
  owners      = [var.account_id]
  most_recent = true

  filter {
    name   = "name"
    values = ["qaas-http-*"]
  }
}

# lc's cannot be updated, unique naming assures most recent version used
resource "aws_launch_configuration" "qaas_server" {
  # omit the "name" attribute to allow Terraform to auto-generate random
  name_prefix                 = "qaas-httpd-server-"
  image_id                    = data.aws_ami.qaas_ami.id
  iam_instance_profile        = aws_iam_instance_profile.qaas_service.name
  instance_type               = "t2.micro"
  associate_public_ip_address = true

  security_groups = [
    aws_security_group.ssh_in.id,
    aws_security_group.http_in.id,
    aws_security_group.http_out.id,
    aws_security_group.ephemeral_out.id,
    aws_security_group.https_out.id,
  ]

  # rolling deployments: create & verify new lc before removing old
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "qaas_service" {
  # interpolate launch configuration name into its name
  # so LC changes always force replacement of the ASG, not just update
  name = "${aws_launch_configuration.qaas_server.name}-asg"
  # Terraform will wait for instances in the new ASG to show up as
  # InService in the ELB before considering the ASG successfully created. 

  min_elb_capacity     = 1
  min_size             = 1
  desired_capacity     = 2
  max_size             = 4
  launch_configuration = aws_launch_configuration.qaas_server.name
  health_check_type    = "ELB"
  vpc_zone_identifier  = var.public_subnets[*].id

  # rolling deployments: create & verify new asg before removing old
  lifecycle {
    create_before_destroy = true
  }

  enabled_metrics = [
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
  target_group_arns = [
    aws_lb_target_group.qaas_target_group.arn,
  ]

  tags = [
    {
      key                 = "qaas"
      value               = "service"
      propagate_at_launch = true
    },
    {
      key                 = "qaas"
      value               = "autoscaling-group"
      propagate_at_launch = false
    },
  ]
}
