variable "codebuild_logging_bucket" {}
variable "codepipeline_artifact_bucket" {}
variable "dynamodb_lock_state_table" {}
variable "faas_iam_role" {}

resource "aws_codebuild_project" "faas_project" {
  name          = "FaaSCodeBuildProject"
  description   = "Terraform codebuild project"
  build_timeout = "5"
  service_role  = var.faas_iam_role.arn

  artifacts {
    type = "CODEPIPELINE"
  }

  cache {
    type     = "S3"
    location = var.codebuild_logging_bucket.bucket
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "travmatth/amazonlinux-golang-dev"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"
    environment_variable {
      name  = "ENV"
      value = "PROD"
      type  = "PLAINTEXT"
    }
  }

  logs_config {
    s3_logs {
      status   = "ENABLED"
      location = "${var.codebuild_logging_bucket.id}/FaaSCodeBuildProject/build-log"
    }
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = file("../../build/cicd/buildspec.yml")
  }

  tags = {
    Terraform = "true"
    FaaS      = "true"
    CodeBuild = "true"
  }
}

output "codebuild_project" {
  value = aws_codebuild_project.faas_project
}