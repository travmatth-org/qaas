output "codedeploy_role_arn" {
	value = aws_iam_role.codedeploy_role.arn
}

output "codedeploy_role_name" {
	value = aws_iam_role.codedeploy_role.name
}