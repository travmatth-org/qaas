output "tf_state_bucket" {
	value = aws_s3_bucket.tf_state_bucket
}

output "codebuild_logging_bucket" {
	value = aws_s3_bucket.codebuild_logging_bucket
}

output "codepipeline_artifact_bucket" {
	value = aws_s3_bucket.codepipeline_artifact_bucket
}

output "dynamodb_lock_state_table" {
	value = aws_dynamodb_table.terraform_lock_state_dynamodb
}

output "tf_iam_assumed_policy" {
	value = aws_iam_policy.tf_iam_assumed_policy
}

output "tf_iam_attach_assumed_role_to_permissions_policy" {
	value = aws_iam_role_policy_attachment.tf_iam_attach_assumed_role_to_permissions_policy
}

output "tf_iam_assumed_role" {
	value = aws_iam_role.tf_iam_assumed_role
}

output "codebuild_iam_role" {
	value = aws_iam_role.codebuild_iam_role
}