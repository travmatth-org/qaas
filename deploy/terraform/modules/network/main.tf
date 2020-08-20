resource "aws_vpc" "local" {
	cidr_block           = "192.168.0.0/16"
	enable_dns_hostnames = true

	tags = {
		Faas = "true"
	}
}

output "public_vpc" {
	value = aws_vpc.local
}

resource "aws_internet_gateway" "faas" {
	vpc_id = aws_vpc.local.id

	tags = {
		FaaS = "gateway"
	}
}

output "internet_gateway" {
	value = aws_internet_gateway.faas
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

output "public_subnet" {
	value = aws_subnet.public_subnet
}

resource "aws_network_acl" "public_acl" {
	vpc_id = aws_vpc.local.id
	subnet_ids = [aws_subnet.public_subnet.id]

	tags = {
		FaaS = "true"
	}
}

resource "aws_network_acl_rule" "http_in" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 100
	egress = false
	protocol = "tcp"
	rule_action = "allow"
	from_port = 80
	to_port = 80
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "http_out" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 100
	egress = true
	protocol = "tcp"
	rule_action = "allow"
	from_port = 80
	to_port = 80
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ssh_in" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 110
	egress = false
	protocol = "tcp"
	rule_action = "allow"
	from_port = 22
	to_port = 22
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ssh_out" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 110
	egress = true
	protocol = "tcp"
	rule_action = "allow"
	from_port = 22
	to_port = 22
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "https_in" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 120
	egress = false
	protocol = "tcp"
	rule_action = "allow"
	from_port = 443
	to_port = 443
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "https_out" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 120
	egress = true
	protocol = "tcp"
	rule_action = "allow"
	from_port = 443
	to_port = 443
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ephemeral_in" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 130
	egress = false
	protocol = "tcp"
	rule_action = "allow"
	from_port = 1024
	to_port = 65535
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ephemeral_out" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 130
	egress = true
	protocol = "tcp"
	rule_action = "allow"
	from_port = 1024
	to_port = 65535
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "deny_in" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 1000
	egress = false
	protocol = "-1"
	rule_action = "deny"
	cidr_block	= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "deny_out" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number = 1000
	egress = true
	protocol = "-1"
	rule_action = "deny"
	cidr_block	= "0.0.0.0/0"
}

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
