resource "aws_instance" "instance" {
	ami = var.ami_id
}