output "codebuild_iam_role_arn" {
  value = aws_iam_role.codebuild_iam_role.arn
}

output "codebuild_iam_role_name" {
  value = aws_iam_role.codebuild_iam_role.name
}