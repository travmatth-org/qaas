data "aws_iam_policy_document" "codebuild_role_policy" {
	statement {
		actions = ["sts:AssumeRole"]

		principals {
			type        = "Service"
			identifiers = ["codebuild.amazonaws.com"]
		}
	}
}

resource "aws_iam_role" "codebuild_role" {
  name               = "FaaSCodeBuildIamRole"
  assume_role_policy = data.aws_iam_policy_document.codebuild_role_policy.json
}

data "aws_iam_policy_document" "codebuild_policy" {
  statement {
    effect    = "Allow"
    actions   = ["s3:*"]
    resources = [
      var.codebuild_logging_bucket.arn,
      "${var.codebuild_logging_bucket.arn}/*",
      var.codepipeline_artifact_bucket.arn,
      "${var.codepipeline_artifact_bucket.arn}/*",
      "arn:aws:s3:::codepipeline-us-west-1*",
      "arn:aws:s3:::codepipeline-us-west-1*/*",
    ]
  }

  statement {
    effect = "Allow"
    resources = ["*"]
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
  }
}

resource "aws_iam_policy" "codebuild_policy" {
  name        = "CodeBuildIamPolicy"
  path        = "/service-role/"
  policy      = data.aws_iam_policy_document.codebuild_policy.json
  description = "Policy used in trust relationship with CodeBuild"
}

resource "aws_iam_policy_attachment" "codebuild_policy_attachment" {
  name        = "CodeBuildPolicyAttachment"
  policy_arn  = aws_iam_policy.codebuild_policy.arn
  roles       = [aws_iam_role.codebuild_role.id]
}

# Create CodeBuild Project for Terraform Plan
resource "aws_codebuild_project" "faas_project" {
  name          = "FaaSCodeBuildProject"
  description   = "Terraform codebuild project"
  build_timeout = "5"
  service_role  = aws_iam_role.codebuild_role.arn

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
    buildspec = file("../../buildspec.yml")
  }

  tags = {
    Terraform = "true"
    FaaS      = "true"
    CodeBuild = "true"
  }
}
