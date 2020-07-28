resource "aws_security_group" "faas_public_http_ssh_sg" {
	name		= "FaaS Public HTTP/SSH"
	vpc_id		= var.public_vpc.id
	description	= "Security group for web that allows web traffic from internet"

	# ingress {
	# 	protocol = "tcp"
	# 	from_port = 22
	# 	to_port = 22
	# 	cidr_blocks = ["0.0.0.0/0"]
	# }

	ingress {
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
		from_port = 0
		to_port = 0
	}

	# allow egress from all ports
	egress {
		from_port = 0
		to_port = 0
		protocol = "-1"
		cidr_blocks = ["0.0.0.0/0"]
	}

	tags = {
		FaaS = "true"
	}
}

data "aws_ami" "amazonlinux2" {
    owners		= ["amazon"]
    most_recent = true

    filter {
        name = "name"
        values = ["amzn2-ami-hvm-2.0.*-x86_64-gp2"]
    }
}

data "aws_iam_policy_document" "faas_service_role" {
	statement {
		effect = "Allow"
		principals {
			type = "Service"
			identifiers	= ["ec2.amazonaws.com"]
		}
		actions = ["sts:AssumeRole"]
	}
}

resource "aws_iam_role" "faas_service_role" {
	name				= "EC2InstanceRole"
	path				= "/"
	assume_role_policy	= data.aws_iam_policy_document.faas_service_role.json

	tags = {
		FaaS = "true"
	}
}

resource "aws_iam_role_policy_attachment" "faas_attachment" {
	role		= aws_iam_role.faas_service_role.name
	policy_arn	= "arn:aws:iam::aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy"
	depends_on	= [aws_iam_role.faas_service_role]
}

resource "aws_iam_instance_profile" "ec2_profile" {
	name		= "faas_service_profile"
	role		= aws_iam_role.faas_service_role.name
	depends_on	= [aws_iam_role.faas_service_role]
}

resource "aws_key_pair" "ec2_key_pair" {
	key_name = "faas_ec2"
	public_key = file("../../protected/faas_ec2.pub")
}

resource "aws_instance" "faas_service" {
	ami						=  data.aws_ami.amazonlinux2.id
	instance_type			= "t2.micro"
	vpc_security_group_ids	= [aws_security_group.faas_public_http_ssh_sg.id]
	subnet_id				= var.public_subnet.id
	iam_instance_profile	= aws_iam_instance_profile.ec2_profile.name
	key_name				= aws_key_pair.ec2_key_pair.key_name
	depends_on				= [var.internet_gateway]
	user_data				= file("../../scripts/s3_user_data.sh")

	tags = {
		faas = "SERVICE"
	}
}

resource "aws_eip" "faas_eip" {
	instance	= aws_instance.faas_service.id
	vpc			= true

	tags = {
		FaaS = "true"
	}
}
