resource "aws_vpc" "local" {
	cidr_block           = "10.0.0.0/16"
	enable_dns_hostnames = true

	tags = {
		faas	= "vpc"
		public	= "true"
	}
}

output "vpc" {
	value = aws_vpc.local
}

resource "aws_internet_gateway" "faas" {
	vpc_id = aws_vpc.local.id

	tags = {
		faas	= "gateway"
		public	= "true"
	}
}

output "internet_gateway" {
	value = aws_internet_gateway.faas
}

data "aws_availability_zones" "available" {}

variable public_subnets {
	description	= "List of public subnets in the VPC"
	type		= list
	default		= [
		"10.0.0.0/18",
		"10.0.64.0/18"
	]
}

variable private_subnets {
	description	= "List of private subnets in the VPC"
	type		= list
	default		= [
		"10.0.128.0/18",
		"10.0.192.0/18"
	]
}


resource "aws_subnet" "public_subnets" {
	count					= 2
	vpc_id                  = aws_vpc.local.id	
	cidr_block				= var.public_subnets[count.index] 
	availability_zone		= data.aws_availability_zones.available.names[count.index]
	map_public_ip_on_launch = true
	depends_on				= [aws_internet_gateway.faas]

	tags = {
		faas	= "subnet"
		public	= "true"
		name	= "public-subnet-az-${count.index}"
	}
}

output "public_subnets" {
	value = aws_subnet.public_subnets
}

resource "aws_route_table" "faas_route_table" {
	vpc_id = aws_vpc.local.id

	route {
		cidr_block = "0.0.0.0/0"
		gateway_id = aws_internet_gateway.faas.id
	}

	tags = {
		faas = "public-route-table"
	}
}

resource "aws_route_table_association" "faas_public" {
	count			= 2
	subnet_id		= aws_subnet.public_subnets[count.index].id
	route_table_id	= aws_route_table.faas_route_table.id
}
