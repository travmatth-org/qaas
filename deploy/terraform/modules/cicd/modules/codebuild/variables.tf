variable "codebuild_logging_bucket" {
  description = "S3 bucket used by CodeBuild for logging"
}

variable "tf_state_bucket" {
  description = "S3 bucket used by CodeBuild for state management"
}

variable "codepipeline_artifact_bucket" {
  description = "s3 bucket used by CodePipeline for artifacts"
}

variable "dynamodb_lock_state_table" {
  description = "DynamoDB table controlling lock state"
}