output "app_name" {
	value = aws_codedeploy_deployment_group.faas_in_place.app_name
}

output "deployment_group_name" {
	value = aws_codedeploy_deployment_group.faas_in_place.deployment_group_name
}