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