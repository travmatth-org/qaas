data "aws_iam_policy_document" "codebuild_role_policy" {
  statement {
    sid     = "QaasCodeBuildTrustRelationships"
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["codebuild.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "codebuild_role" {
  name                  = "QaasCodeBuildIamRole"
  assume_role_policy    = data.aws_iam_policy_document.codebuild_role_policy.json
  force_detach_policies = true
}

data "aws_iam_policy_document" "codebuild_policy" {
  statement {
    sid     = "AllowCodeBuildS3Control"
    effect  = "Allow"
    actions = ["s3:*"]
    resources = [
      var.codebuild_logging_bucket.arn,
      "${var.codebuild_logging_bucket.arn}/*",
      "${var.codepipeline_artifact_bucket.arn}/*",
      "arn:aws:s3:::codepipeline-us-west-1/*",
    ]
  }

  statement {
    sid    = "AllowCodeBuildLoggingAccess"
    effect = "Allow"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["*"]
  }

  statement {
    sid     = "AllowCodeBuildSSMParameterAccess"
    effect  = "Allow"
    actions = ["ssm:GetParameters"]
    resources = [
      "arn:aws:ssm:us-west-1:${var.account_id}:parameter/qaas/packer-role-arn"
    ]
  }

  statement {
    sid    = "AllowCodeBuildToAssumePackerRole"
    effect = "Allow"
    actions = [
      "sts:AssumeRole",
    ]
    resources = [
      aws_iam_role.packer_sts_delegate.arn
    ]
  }
}

resource "aws_iam_policy" "codebuild_policy" {
  name   = "QaasCodeBuildPolicy"
  path   = "/service-role/"
  policy = data.aws_iam_policy_document.codebuild_policy.json
}

resource "aws_iam_policy_attachment" "codebuild_policy_attachment" {
  name       = "QaasCodeBuildPolicyAttachment"
  policy_arn = aws_iam_policy.codebuild_policy.arn
  roles      = [aws_iam_role.codebuild_role.id]
}