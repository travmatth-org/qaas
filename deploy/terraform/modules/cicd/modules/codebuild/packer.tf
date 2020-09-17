data "aws_iam_policy_document" "packer_sts_delegate" {
	statement {
		actions = ["sts:AssumeRole"]

		principals {
			type		= "AWS"
			identifiers	= ["arn:aws:iam::${var.account_id}:user/${var.user_name}"]
		}
	}
}

resource "aws_iam_role" "packer_sts_delegate" {
	name				= "packer-sts-delegate"
	assume_role_policy	= data.aws_iam_policy_document.packer_sts_delegate.json
}

data "aws_iam_policy_document" "codebuild_packer_policy_document" {
	statement {
		sid		= "AllowPackerAmiCreation"
		effect	= "Allow"
		actions	= [
			"ec2:AttachVolume",
			"ec2:AuthorizeSecurityGroupIngress",
			"ec2:CopyImage",
			"ec2:CreateImage",
			"ec2:CreateKeypair",
			"ec2:CreateSecurityGroup",
			"ec2:CreateSnapshot",
			"ec2:CreateTags",
			"ec2:CreateVolume",
			"ec2:DeleteKeyPair",
			"ec2:DeleteSecurityGroup",
			"ec2:DeleteSnapshot",
			"ec2:DeleteVolume",
			"ec2:DeregisterImage",
			"ec2:DescribeImageAttribute",
			"ec2:DescribeImages",
			"ec2:DescribeInstances",
			"ec2:DescribeInstanceStatus",
			"ec2:DescribeRegions",
			"ec2:DescribeSecurityGroups",
			"ec2:DescribeSnapshots",
			"ec2:DescribeSubnets",
			"ec2:DescribeTags",
			"ec2:DescribeVolumes",
			"ec2:DescribeVpcs",
			"ec2:DetachVolume",
			"ec2:ModifyImageAttribute",
			"ec2:ModifyInstanceAttribute",
			"ec2:ModifySnapshotAttribute",
			"ec2:RegisterImage",
			"ec2:RunInstances",
			"ec2:StopInstances",
			"ec2:TerminateInstances",
			"iam:PassRole",
			"iam:GetInstanceProfile"
		]
		resources = ["*"]
	}
}

resource "aws_iam_policy" "codebuild_packer_policy" {
	name		= "codebuild-packer-policy"
	description	= "Allow packer to build AMI within codebuild"
	policy		= data.aws_iam_policy_document.codebuild_packer_policy_document.json
}

resource "aws_iam_policy_attachment" "codebuild_packer_policy" {
	name		= "packer-codebuild-policy"
	policy_arn	= aws_iam_policy.codebuild_packer_policy.arn
	roles		= [aws_iam_role.packer_sts_delegate.id]
}

resource "aws_ssm_parameter" "packer_arn_secret" {
	name		= "/faas/packer_role_arn"
	description	= "ARN of the role for packer to build AMI within codebuild"
	type		= "SecureString"
	value		= aws_iam_role.packer_sts_delegate.arn
	tags		= {
		faas	= "secret-param"
	}
}