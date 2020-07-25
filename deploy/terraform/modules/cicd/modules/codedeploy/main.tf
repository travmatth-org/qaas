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

resource "aws_iam_role_policy_attachment" "codedeploy_attach" {
	policy_arn	= "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole"
	role		= aws_iam_role.codedeploy_role.name
	depends_on	= [aws_iam_role.codedeploy_role]
}

resource "aws_codedeploy_deployment_group" "faas_in_place" {
	app_name 			  = aws_codedeploy_app.faas.name
	deployment_group_name = "${aws_codedeploy_app.faas.name}-deployment-group"
	service_role_arn 	  = aws_iam_role.codedeploy_role.arn

	ec2_tag_set {
		ec2_tag_filter {
			type  = "KEY_AND_VALUE"	
			key   = "FaaS"
			value = "Service"
		}
	}

	deployment_style {
		deployment_type = "IN_PLACE"
	}
}
