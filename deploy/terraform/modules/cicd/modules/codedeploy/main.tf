resource "aws_codedeploy_app" "faas" {
	name = "faas"
}

resource "aws_iam_role" "codedeploy_role" {
	name = "FaasCodeDeployRole"
	assume_role_policy =<<EOF
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Sid": "",
			"Effect": "Allow",
			"Principal": {
				"Service": "codedeploy.amazon.aws.com"
			},
			"Action": "sts:AssumeRole"
		}
	]
}
EOF
}

resource "aws_iam_role_policy_attachment" "AWSCodeDeployRole" {
	policy_arn = "arn:aws:iam:policy/service-role/AWSCodeDeployRole"
	role = aws_iam_role.codedeploy_role.name
}

resource "aws_codedeploy_deployment_group" "deploy" {
	app_name 			  = aws_codedeploy_app.faas.name
	deployment_group_name = "faas-ec2-group"
	service_role_arn 	  = aws_iam_role.codedeploy_role.arn

	ec2_tag_set {
		ec2_tag_filter {
			type  = "KEY_AND_VALUE"	
			key   = "FaaS"
			value = "Service"
		}
	}
}
