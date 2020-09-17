variable "codepipeline_artifact_bucket" {}
variable "user" {}

data "aws_iam_policy_document" "assume_role" {
	statement {
		effect  = "Allow"
		actions = ["sts:AssumeRole"]

		principals {
			type = "Service"
			identifiers = [
				"ec2.amazonaws.com",
				"codedeploy.amazonaws.com",
				"s3.amazonaws.com",
			]
		}
	}
}

resource "aws_iam_role" "faas" {
  name					= "faas_role"
  assume_role_policy	= data.aws_iam_policy_document.assume_role.json
}

output "faas_iam_role" {
	value = aws_iam_role.faas
}

data "aws_iam_policy_document" "faas_cicd_policy" {
	statement {
		sid		= "AllowAgentMetadataQueries"
		effect	= "Allow"
		actions = [
			"ec2:DescribeTags"
		]
		resources = ["*"]
	}

	statement {
		sid		= "AllowCodeDeployDeploymnets"
		effect	= "Allow"
		actions	= [
			"s3:Get*",
			"s3:List*",
		]
		resources = [
			"${var.codepipeline_artifact_bucket.arn}/*",
			"arn:aws:s3:::aws-codedeploy-us-west-1/*"
		]
	}
	
	# https://github.com/SummitRoute/aws_managed_policies/blob/master/policies/CloudWatchAgentServerPolicy
	statement {
		sid		= "AllowCloudWatchLogging"
		effect  = "Allow"
		actions = [
			"cloudwatch:PutMetricData",
			"logs:PutLogEvents",
			"logs:DescribeLogStreams",
			"logs:DescribeLogGroups",
			"logs:CreateLogStream",
		]
		resources = ["*"]
	}

	# https://docs.aws.amazon.com/xray/latest/devguide/security_iam_id-based-policy-examples.html
	statement {
		sid		= "AllowXRayTraces"
		effect  = "Allow"
		actions = [
			"xray:PutTraceSegments",
			"xray:PutTelemetryRecords",
			"xray:GetSamplingRules",
			"xray:GetSamplingTargets",
			"xray:GetSamplingStatisticSummaries"
		]
		resources = ["*"]
	}

	# allow ssm to manage ssh to instances
	statement {
		sid		= "AllowSSMSSHAccess"
		effect	= "Allow"
		actions	= [
			"ssmmessages:CreateControlChannel",
			"ssmmessages:CreateDataChannel",
			"ssmmessages:OpenControlChannel",
			"ssmmessages:OpenDataChannel"
		]
		resources = ["*"]
	}
}

resource "aws_iam_policy" "cicd" {
  name			= "faas_role_policy"
  description	= "A policy for CodeDeploy"
  policy		= data.aws_iam_policy_document.faas_cicd_policy.json
  depends_on	= [aws_iam_role.faas]
}

resource "aws_iam_role_policy_attachment" "attach_custom_policy" {
  depends_on = [aws_iam_policy.cicd]
  role       = aws_iam_role.faas.name
  policy_arn = aws_iam_policy.cicd.arn
}

resource "aws_iam_instance_profile" "faas_service" {
  name = "faas_service"
  path = "/"
  role = aws_iam_role.faas.name
}

resource "aws_iam_group" "faas_ssh_group" {
	name	= "faas-ssh-group"
}

resource "aws_iam_group_membership" "faas_ssh_group" {
	name	= "ssh-group-members"
	users	= [var.user]
	group	= aws_iam_group.faas_ssh_group.id
}

data "aws_iam_policy_document" "ssh_policy" {
	statement {
		sid			= "AllowUserSSHAccessFaasInstance"
		effect		= "Allow"
		actions		= ["ec2-instance-connect:SendSSHPublicKey"]
		resources	= ["arn:aws:ec2:us-west-1:${var.account_id}:instance/*"]
		condition {
			test		= "StringEquals"
			variable	= "ec2:osuser"
			values		= ["ec2-user"]
		}
		condition {
			test		= "StringEquals"
			variable	= "aws:ResourceTag/faas"
			values		= ["service"]
		}
	}

	# EC2 Instance Connect CLI wrapper calls this action
	statement {
		effect    = "Allow"
		actions   = ["ec2:DescribeInstances"]
		resources = ["*"]
	}
}

resource "aws_iam_group_policy" "ssh_policy" {
	name	= "faas_ssh-policy"
	group	= aws_iam_group.faas_ssh_group.id

	policy	= data.aws_iam_policy_document.ssh_policy.json
}
