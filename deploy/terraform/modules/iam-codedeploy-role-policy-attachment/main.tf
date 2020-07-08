resource "aws_iam_role_policy_attachment" "AWSCodeDeployRole" {
	policy_arn = "arn:aws:iam:policy/service-role/AWSCodeDeployRole"
	role = "${var.codedeploy_role_name}"
}