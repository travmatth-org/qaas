variable "s3_tfstate_bucket" {
  description = "Name of the S3 bucket used for Terraform state storage"
}

variable "s3_logging_bucket_name" {
  description = "Name of S3 bucket to use for access logging"
}

variable "tf_codepipeline_artifact_bucket_arn" {
  description = "Codepipeline artifact bucket ARN"
}