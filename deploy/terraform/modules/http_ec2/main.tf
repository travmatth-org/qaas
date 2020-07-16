resource "aws_vpc" "faas_vpc" {
	cidr_block           = "10.0.0.0/16"
	enable_dns_hostnames = true

	tags = {
		Faas = "true"
	}
}

resource "aws_subnet" "faas_subnet" {
	vpc_id                  = aws_vpc.faas_vpc.id
	cidr_block              = "10.0.1.0/24"
	map_public_ip_on_launch = true 
	availability_zone       = "us-west-1"

	tags = {
		FaaS = "true"
	}
}

resource "aws_security_group" "faas_public_http_ssh_sg" {
	vpc_id = aws_vpc.faas_vpc.id
	name = "FaaS Public HTTP/SSH"
	description = "Allow incoming SSH traffic, all egress"

	ingress {
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
		from_port = 22
		to_port = 22
	}

	ingress {
		protocol = "tcp"
		cidr_blocks = ["0.0.0.0/0"]
		from_port = 80
		to_port = 80
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

resource "aws_network_acl" "faas_public_subnet_acl" {
	vpc_id = aws_vpc.faas_vpc.id
	subnet_ids = [aws_subnet.faas_subnet.id]

	# ssh ingress, port 22
	ingress {
		protocol = "tcp"
		rule_no = 100
		action = "allow"
		cidr_block = "0.0.0.0/0"
		from_port = 22
		to_port = 22
	}

	# http ingress, port 80
	ingress {
		protocol = "tcp"
		rule_no = 100
		action = "allow"
		cidr_block = "0.0.0.0/0"
		from_port = 80
		to_port = 80
	}

	# egress, ephemeral ports
	egress {
		protocol = "tcp"
		rule_no = 300
		action = "allow"
		cidr_block = "0.0.0.0/0"
		from_port = 1024
		to_port = 65535
	}

	tags = {
		FaaS = "true"
	}
}

resource "aws_internet_gateway" "faas_gateway" {
	vpc_id = aws_vpc.faas_vpc.id

	tags = {
		FaaS = "true"
	}
}

resource "aws_route_table" "faas_route_table" {
	vpc_id = aws_vpc.faas_vpc.id

	route {
		cidr_block = "0.0.0.0/0"
		gateway_id = aws_internet_gateway.faas_gateway.id
	}

	tags = {
		FaaS = "true"
	}
}

resource "aws_route" "faas_route" {
	route_table_id = aws_route_table.faas_route_table.id
	destination_cidr_block = "0.0.0.0/0"
	gateway_id = aws_internet_gateway.faas_gateway.id
}

resource "aws_route_table_association" "faas_assoc" {
	subnet_id = aws_subnet.faas_subnet.id
	route_table_id = aws_route.faas_route.id
}

data "aws_ami" "amazonlinux2" {
    owners = ["amazon"]
    most_recent = true
    filter {
        name = "name"
        values = ["amzn2-ami-hvm-2.0.*-x86_64-gp2"]
    }
}


resource "aws_iam_role" "faas_service_role" {
	assume_role_policy = <<-EOF
		{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}]
		}
		EOF

	tags = {
		FaaS = "true"
	}
}

resource "aws_iam_role_policy_attachment" "faas_attachment" {
	role = aws_iam_role.faas_service_role.name
	policy_arn = "arn:aws:iam:aws:policy/service-role/AmazonEC2RoleforAWSCodeDeploy"
}

resource "aws_iam_instance_profile" "ec2_profile" {
	name = "faas_service_profile"
	role = aws_iam_role.faas_service_role.name
}

resource "aws_key_pair" "ec2_key_pair" {
	key_name = "faas_ec2"
	public_key = file("../../protected/faas_ec2.pub")
}

resource "aws_instance" "faas_service" {
	ami =  data.aws_ami.amazonlinux2.id
	depends_on = [aws_internet_gateway.faas_gateway]
	instance_type = "t2.micro"
	security_groups = [aws_security_group.faas_public_http_ssh_sg.id]
	subnet_id = aws_subnet.faas_subnet.id
	iam_instance_profile = aws_iam_instance_profile.ec2_profile.name
	key_name = aws_key_pair.ec2_key_pair.key_name
	user_data = <<-EOF
		#!/bin/bash
		sudo yum -y update && \
			install ruby wget
		cd /home/ec2-user
		wget https://aws-codedeploy-us-west-1.s3.amazonaws.com/latest/install
		chmod +x ./install
		sudo ./intall auto
		sudo service codedeploy-agent status
		EOF

	tags = {
		FaaS = "Service"
	}
}

resource "aws_eip" "faas_eip" {
	instance = aws_instance.faas_service.id
	vpc = true

	tags = {
		FaaS = "true"
	}
}