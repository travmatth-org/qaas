output "logging_bucket_id" {
  value = aws_s3_bucket.logging_bucket.id
}

output "logging_bucket" {
  value = aws_s3_bucket.logging_bucket.bucket
}

output "state_bucket_id" {
  value = aws_s3_bucket.state_bucket.id
}

output "state_bucket" {
  value = aws_s3_bucket.state_bucket.bucket
}

output "codepipeline_artifact_bucket_id" {
  value = aws_s3_bucket.codepipeline_artifact_bucket.id
}

output "codepipeline_artifact_bucket" {
  value = aws_s3_bucket.codepipeline_artifact_bucket.bucket
}