resource "aws_iam_role_policy" "codebuild_iam_role_policy" {
  name = "CodeBuildIamRolePolicy"
  role = var.codebuild_iam_role.name
  policy = <<-POLICY
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Resource": [
            "*"
          ],
          "Action": [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "s3:*"
          ],
          "Resource": [
            "${var.codebuild_logging_bucket.arn}",
            "${var.codebuild_logging_bucket.arn}/*",
            "${var.tf_state_bucket.arn}",
            "${var.tf_state_bucket.arn}/*",
            "arn:aws:s3:::codepipeline-us-west-1*",
            "arn:aws:s3:::codepipeline-us-west-1*/*",
            "${var.codepipeline_artifact_bucket.arn}",
            "${var.codepipeline_artifact_bucket.arn}/*"
          ]
        },
        {
          "Effect": "Allow",
          "Action": [
            "dynamodb:*"
          ],
          "Resource": "${var.dynamodb_lock_state_table.arn}"
        },
        {
          "Effect": "Allow",
          "Action": [
            "iam:Get*",
            "iam:List*"
          ],
          "Resource": "${var.codebuild_iam_role.arn}"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Resource": "${var.codebuild_iam_role.arn}"
        }
      ]
    }
    POLICY
}

# Create CodeBuild Project for Terraform Plan
resource "aws_codebuild_project" "faas_project" {
  name          = "FaaSCodeBuildProject"
  description   = "Terraform codebuild project"
  build_timeout = "5"
  service_role  = var.codebuild_iam_role.arn

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
    # environment_variable {
    # }
  }

  logs_config {
    s3_logs {
      status   = "ENABLED"
      location = "${var.codebuild_logging_bucket.id}/FaaSCodeBuildProject/build-log"
    }
  }

  source {
    type      = "CODEPIPELINE"
    buildspec = var.buildspec_yml
  }

  tags = {
    Terraform = "true"
    FaaS      = "true"
    CodeBuild = "true"
  }
}
