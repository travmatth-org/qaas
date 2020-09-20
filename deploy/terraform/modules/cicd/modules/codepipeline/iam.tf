data "aws_iam_policy_document" "codepipeline_role_policy" {
  statement {
    sid    = "QaasCodePipelineTrustRelationships"
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["codepipeline.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "codepipeline_role" {
  name               = "CodePipelineIamRole"
  assume_role_policy = data.aws_iam_policy_document.codepipeline_role_policy.json

  tags = {
    qaas = "true"
  }
}

variable "account_id" {}

data "aws_iam_policy_document" "codepipeline_policy" {
  statement {
    sid    = "QaasCodePipelineS3Object"
    effect = "Allow"
    actions = [
      "s3:GetObject*",
      "s3:PutObject",
      "s3:PutObjectAcl",
    ]
    resources = [
      "${var.codebuild_logging_bucket.arn}/*",
      "${var.codepipeline_artifact_bucket.arn}/*",
      "arn:aws:s3:::codepipeline-us-west-1/*",
    ]
  }

  statement {
    sid    = "QaasCodepipelineAllowCodeBuild"
    effect = "Allow"
    actions = [
      "codebuild:BatchGetBuilds",
      "codebuild:StartBuild",
    ]
    resources = [var.codebuild_project.id]
  }

  statement {
    sid    = "QaasCodepipelineAllowCodeDeployRevision"
    effect = "Allow"
    actions = [
      "codedeploy:RegisterApplicationRevision",
      "codedeploy:GetApplicationRevision"
    ]
    resources = [
      "arn:aws:codedeploy:us-west-1:${var.account_id}:application:${var.codedeploy.app.name}"
    ]
  }

  statement {
    sid    = "QaasCodepipelineAllowCodeDeployApplication"
    effect = "Allow"
    actions = [
      "codedeploy:CreateDeployment",
      "codedeploy:GetDeployment",
    ]
    resources = [
      "arn:aws:codedeploy:us-west-1:${var.account_id}:deploymentgroup:${var.codedeploy.app.name}/${var.codedeploy.deployment_group.deployment_group_name}"
    ]
  }

  statement {
    sid    = "QaasCodepipelineAllowCodeDeployDeployment"
    effect = "Allow"
    actions = [
      "codedeploy:GetDeploymentConfig"
    ]
    resources = [
      "arn:aws:codedeploy:us-west-1:${var.account_id}:deploymentconfig:${var.codedeploy.deployment_group.deployment_config_name}",
    ]
  }
}

resource "aws_iam_role_policy" "codepipeline_policy" {
  name   = "QaasCodePipelinePolicy"
  role   = aws_iam_role.codepipeline_role.id
  policy = data.aws_iam_policy_document.codepipeline_policy.json
}
