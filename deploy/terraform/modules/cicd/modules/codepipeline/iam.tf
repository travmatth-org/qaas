data "aws_iam_policy_document" "codepipeline_role_policy" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["codepipeline.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "codepipeline_role" {
  name               = "TerraformCodePipelineIamRole"
  assume_role_policy = data.aws_iam_policy_document.codepipeline_role_policy.json

  tags               = {
    FaaS = "true"
  }
}

data "aws_iam_policy_document" "codepipeline_policy" {
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
    effect    = "Allow"
    actions   = [
      "codebuild:BatchGetBuilds",
      "codebuild:StartBuild",
    ]
    resources = ["*"]
  }

  statement {
    effect    = "Allow"
    actions   = ["codedeploy:*"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "codepipeline_policy" {
  name    = "FaasCodePipelinePolicy"
  role    = aws_iam_role.codepipeline_role.id
  policy  = data.aws_iam_policy_document.codepipeline_policy.json
}