output "public_vpc" {
	value = aws_vpc.public_vpc
}

output "public_subnet" {
	value = aws_subnet.public_subnet
}

output "internet_gateway" {
	value = aws_internet_gateway.faas_gateway
}