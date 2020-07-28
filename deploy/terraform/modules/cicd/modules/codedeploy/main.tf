resource "aws_codedeploy_app" "faas" {
	name = "faas"
}

data "aws_iam_policy_document" "codedeploy" {
	statement {
		actions = ["sts:AssumeRole"]

		principals {
			type		= "Service"
			identifiers = ["codedeploy.amazonaws.com"]
		}
	}
}

resource "aws_iam_role" "codedeploy_role" {
	name				= "faas_codedeploy_role"
	assume_role_policy	= data.aws_iam_policy_document.codedeploy.json
	path				= "/"
	description			= "Allows CodeDeploy to call AWS services"
}

resource "aws_iam_role_policy" "codedeploy" {
	name = "FaasRolePolicy"
	role = aws_iam_role.codedeploy_role.name
	policy = <<-POLICY
		{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": [
						"s3:Get*",
						"s3:List*"
					],
					"Effect": "Allow",
					"Resource": [
						"${var.codepipeline_artifact_bucket.arn}/*",
						"arn:aws:s3:::aws-codedeploy-us-west-1/*"
					]
				}
			]
		}
		POLICY
}

resource "aws_iam_role_policy_attachment" "codedeploy_attach" {
	policy_arn	= "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole"
	role		= aws_iam_role.codedeploy_role.name
}

resource "aws_codedeploy_deployment_group" "faas_in_place" {
	app_name 			  = aws_codedeploy_app.faas.name
	deployment_group_name = "${aws_codedeploy_app.faas.name}-deployment-group"
	service_role_arn 	  = aws_iam_role.codedeploy_role.arn

	ec2_tag_set {
		ec2_tag_filter {
			type  = "KEY_AND_VALUE"	
			key   = "faas"
			value = "SERVICE"
		}
	}

	deployment_style {
		deployment_type = "IN_PLACE"
	}
}
