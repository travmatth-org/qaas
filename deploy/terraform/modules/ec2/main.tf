variable "public_subnet" {}
variable "internet_gateway" {}
variable "codepipeline_artifact_bucket" {}

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
	vpc_security_group_ids	= [
		aws_security_group.http_in.id,
		aws_security_group.http_out.id,
		aws_security_group.ssh_in.id,
		aws_security_group.ephemeral_out.id,
		aws_security_group.https_out.id,
	]
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