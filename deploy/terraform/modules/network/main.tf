resource "aws_vpc" "public_vpc" {
	cidr_block           = "10.0.0.0/16"
	enable_dns_hostnames = true

	tags = {
		Faas = "true"
	}
}

resource "aws_internet_gateway" "faas_gateway" {
	vpc_id = aws_vpc.public_vpc.id

	tags = {
		FaaS = "true"
	}
}

resource "aws_subnet" "public_subnet" {
	vpc_id                  = aws_vpc.public_vpc.id
	cidr_block              = "10.0.1.0/24"
	map_public_ip_on_launch = true 
	availability_zone       = "us-west-1b"
	depends_on				= [aws_internet_gateway.faas_gateway]

	tags = {
		FaaS = "true"
	}
}

resource "aws_network_acl" "public_subnet_acl" {
	vpc_id = aws_vpc.public_vpc.id
	subnet_ids = [aws_subnet.public_subnet.id]

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
		rule_no = 101
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

resource "aws_route_table" "faas_route_table" {
	vpc_id = aws_vpc.public_vpc.id

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
	subnet_id = aws_subnet.public_subnet.id
	route_table_id = aws_route_table.faas_route_table.id
}
