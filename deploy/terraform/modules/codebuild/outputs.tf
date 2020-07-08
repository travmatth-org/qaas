output "project_name" {
  value = var.project_name
}

output "codebuild_project_arn" {
  value = aws_codebuild_project.faas_project.arn
}