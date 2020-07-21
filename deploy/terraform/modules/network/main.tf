resource "aws_vpc" "local" {
	cidr_block           = "192.168.0.0/16"
	enable_dns_hostnames = true

	tags = {
		Faas = "true"
	}
}

resource "aws_internet_gateway" "faas" {
	vpc_id = aws_vpc.local.id

	tags = {
		FaaS = "gateway"
	}
}

resource "aws_subnet" "public_subnet" {
	vpc_id                  = aws_vpc.local.id
	cidr_block              = "192.168.0.0/24"
	map_public_ip_on_launch = true 
	availability_zone       = "us-west-1b"
	depends_on				= [aws_internet_gateway.faas]

	tags = {
		FaaS = "public-subnet"
	}
}

# resource "aws_network_acl" "public_subnet_acl" {
# 	vpc_id = aws_vpc.local.id
# 	subnet_ids = [aws_subnet.public_subnet.id]

	# ingress all
	# ingress {
	# 	protocol	= "tcp"
	# 	rule_no		= 100
	# 	action		= "allow"
	# 	cidr_block	= "0.0.0.0/0"
	# 	from_port	= 0
	# 	to_port		= 65535
	# }

	# egress all
	# ingress {
	# 	protocol	= "tcp"
	# 	rule_no		= 200
	# 	action		= "allow"
	# 	cidr_block	= "0.0.0.0/0"
	# 	from_port	= 0
	# 	to_port		= 65535
	# }


	# ssh ingress, port 22
	# ingress {
	# 	protocol = "tcp"
	# 	rule_no = 100
	# 	action = "allow"
	# 	cidr_block = "0.0.0.0/0"
	# 	from_port = 22
	# 	to_port = 22
	# }

	# http ingress, port 80
	# ingress {
	# 	protocol = "tcp"
	# 	rule_no = 101
	# 	action = "allow"
	# 	cidr_block = "0.0.0.0/0"
	# 	from_port = 80
	# 	to_port = 80
	# }

	# egress, ephemeral ports
	# egress {
	# 	protocol = "tcp"
	# 	rule_no = 300
	# 	action = "allow"
	# 	cidr_block = "0.0.0.0/0"
	# 	from_port = 1024
	# 	to_port = 65535
	# }

# 	tags = {
# 		FaaS = "public-acl"
# 	}
# }

resource "aws_route_table" "faas_route_table" {
	vpc_id = aws_vpc.local.id

	route {
		cidr_block = "0.0.0.0/0"
		gateway_id = aws_internet_gateway.faas.id
	}

	tags = {
		FaaS = "public-route-table"
	}
}

# resource "aws_route" "faas_route" {
# 	route_table_id = aws_route_table.faas_route_table.id
# 	destination_cidr_block = "0.0.0.0/0"
# 	gateway_id = aws_internet_gateway.faas.id
# }

resource "aws_route_table_association" "faas_public" {
	subnet_id = aws_subnet.public_subnet.id
	route_table_id = aws_route_table.faas_route_table.id
}
