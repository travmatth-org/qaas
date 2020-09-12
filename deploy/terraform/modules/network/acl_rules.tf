resource "aws_network_acl" "public_acl" {
	vpc_id		= aws_vpc.local.id
	subnet_ids	= aws_subnet.public_subnets[*].id

	tags = {
		faas = "true"
	}
}

resource "aws_network_acl_rule" "http_in" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 100
	egress			= false
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 80
	to_port			= 80
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "http_out" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 100
	egress			= true
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 80
	to_port			= 80
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ssh_in" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 110
	egress			= false
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 22
	to_port			= 22
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ssh_out" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 110
	egress			= true
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 22
	to_port			= 22
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "https_in" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 120
	egress			= false
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 443
	to_port			= 443
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "https_out" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 120
	egress			= true
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 443
	to_port			= 443
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ephemeral_in" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 130
	egress			= false
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 1024
	to_port			= 65535
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "ephemeral_out" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 130
	egress			= true
	protocol		= "tcp"
	rule_action		= "allow"
	from_port		= 1024
	to_port			= 65535
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "deny_in" {
	network_acl_id	= aws_network_acl.public_acl.id
	rule_number		= 1000
	egress			= false
	protocol		= "-1"
	rule_action		= "deny"
	cidr_block		= "0.0.0.0/0"
}

resource "aws_network_acl_rule" "deny_out" {
	network_acl_id = aws_network_acl.public_acl.id
	rule_number	= 1000
	egress		= true
	protocol	= "-1"
	rule_action	= "deny"
	cidr_block	= "0.0.0.0/0"
}
