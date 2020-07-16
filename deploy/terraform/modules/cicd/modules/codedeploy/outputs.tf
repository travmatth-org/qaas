output "app_name" {
	value = aws_codedeploy_deployment_group.deploy.app_name
}

output "deployment_group_name" {
	value = aws_codedeploy_deployment_group.deploy.deployment_group_name
}