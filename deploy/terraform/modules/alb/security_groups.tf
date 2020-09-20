resource "aws_security_group" "http_in" {
  name        = "qaas-http-in"
  vpc_id      = var.vpc.id
  description = "allow http ingress traffic"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    qaas = "true"
  }
}

resource "aws_security_group" "http_out" {
  name        = "qaas-http-out"
  vpc_id      = var.vpc.id
  description = "allow http egress traffic"

  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    qaas = "true"
  }
}

resource "aws_security_group" "ssh_in" {
  name        = "qaas-ssh-in"
  vpc_id      = var.vpc.id
  description = "allow ssh ingress traffic"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    qaas = "true"
  }
}

resource "aws_security_group" "ephemeral_out" {
  name        = "qaas-ephemeral-out"
  vpc_id      = var.vpc.id
  description = "allow egress traffic to ephemeral ports"

  egress {
    from_port   = 1024
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    qaas = "true"
  }
}


resource "aws_security_group" "https_out" {
  name        = "qaas-https-out"
  vpc_id      = var.vpc.id
  description = "allow https egress traffic"

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    qaas = "true"
  }
}