output "codepipeline_role" {
	value = aws_iam_role.codepipeline_role
}

output "codepipeline" {
	value = aws_codepipeline.codepipeline
}