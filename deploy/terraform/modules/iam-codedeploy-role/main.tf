resource "aws_iam_role" "codedeploy_role" {
	name = "faas_codedeploy_role"

	assume_role_policy = <<EOF
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