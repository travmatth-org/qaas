variable "public_vpc" {}
variable "public_subnet" {}
variable "internet_gateway" {}
variable "codepipeline_artifact_bucket" {}

resource "aws_security_group" "faas_public_http_ssh_sg" {
	name		= "FaaS Public HTTP/SSH"
	vpc_id		= var.public_vpc.id
	description	= "Security group for web that allows web traffic from internet"

	ingress {
		from_port = 80
		to_port = 80
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
	}

	ingress {
		from_port = 22
		to_port = 22
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
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

resource "aws_key_pair" "ec2_key_pair" {
	key_name = "faas_ec2"
	public_key = file("../../protected/faas_ec2.pub")
}

resource "aws_instance" "faas_service" {
	ami						=  data.aws_ami.amazonlinux2.id
	instance_type			= "t2.micro"
	vpc_security_group_ids	= [aws_security_group.faas_public_http_ssh_sg.id]
	subnet_id				= var.public_subnet.id
	iam_instance_profile	= aws_iam_instance_profile.faas_service.name
	key_name				= aws_key_pair.ec2_key_pair.key_name
	depends_on				= [var.internet_gateway]
	user_data				= file("../../scripts/ec2/user_data.sh")

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

output "faas_eip" {
	value	= aws_eip.faas_eip.public_ip
}