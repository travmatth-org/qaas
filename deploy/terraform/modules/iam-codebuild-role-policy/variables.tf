variable "codebuild_iam_role_policy_name" {
  description = "Name for IAM policy used by CodeBuild"
}

variable "codebuild_iam_role_name" {
  description = "Name for IAM role used by CodeBuild"
}

variable "logging_bucket_arn" {
  description = "ARN for s3 bucket used by CodeBuild for logging"
}

variable "state_bucket_arn" {
  description = "ARN for s3 bucket used by CodeBuild for state management"
}

variable "codepipeline_artifact_bucket_arn" {
  description = "ARN for s3 bucket used by CodePipeline for artifacts"
}

variable "dynamodb_lock_state_table_arn" {
  description = "ARN of the dynamodb table controlling lock state"
}

variable "codebuild_iam_role_arn" {
  description = "ARN of the CodeBuild Role"
}