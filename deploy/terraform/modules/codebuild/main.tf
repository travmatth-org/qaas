# Create CodeBuild Project for Terraform Plan
resource "aws_codebuild_project" "faas_project" {
  name          = var.project_name
  description   = "Terraform codebuild project"
  build_timeout = "5"
  service_role  = var.codebuild_iam_role_arn

  artifacts {
    type = "CODEPIPELINE"
  }

  cache {
    type     = "S3"
    location = var.s3_logging_bucket
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "travmatth/amazonlinux-golang-dev"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    environment_variable {
      name  = "TERRAFORM_VERSION"
      value = "0.12.16"
    }
  }

  logs_config {
    s3_logs {
      status   = "ENABLED"
      location = "${var.s3_logging_bucket_id}/${var.project_name}/build-log"
    }
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = "/build/ci/buildspec.yml"
  }

  tags = {
    Terraform = "true"
    FaaS      = "true"
    CodeBuild = "true"
  }
}
