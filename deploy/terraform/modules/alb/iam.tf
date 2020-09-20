data "aws_iam_policy_document" "assume_role" {
  statement {
    sid     = "TrustRelationships"
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type = "Service"
      identifiers = [
        "ec2.amazonaws.com",
        "codedeploy.amazonaws.com",
        "s3.amazonaws.com",
      ]
    }
  }
}

resource "aws_iam_role" "qaas" {
  name               = "QaasRole"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

output "qaas_iam_role" {
  value = aws_iam_role.qaas
}

data "aws_iam_policy_document" "qaas_cicd_policy" {
  statement {
    sid    = "AllowAgentMetadataQueries"
    effect = "Allow"
    actions = [
      "ec2:DescribeTags"
    ]
    resources = ["*"]
  }

  statement {
    sid    = "AllowCodeDeployDeploymnets"
    effect = "Allow"
    actions = [
      "s3:Get*",
      "s3:List*",
    ]
    resources = [
      "${var.codepipeline_artifact_bucket.arn}/*",
      "arn:aws:s3:::aws-codedeploy-us-west-1/*"
    ]
  }

  # https://github.com/SummitRoute/aws_managed_policies/blob/master/policies/CloudWatchAgentServerPolicy
  statement {
    sid    = "AllowCloudWatchLogging"
    effect = "Allow"
    actions = [
      "cloudwatch:PutMetricData",
      "logs:PutLogEvents",
      "logs:DescribeLogStreams",
      "logs:DescribeLogGroups",
      "logs:CreateLogStream",
    ]
    resources = ["*"]
  }

  # https://docs.aws.amazon.com/xray/latest/devguide/security_iam_id-based-policy-examples.html
  statement {
    sid    = "AllowXRayTraces"
    effect = "Allow"
    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords",
      "xray:GetSamplingRules",
      "xray:GetSamplingTargets",
      "xray:GetSamplingStatisticSummaries"
    ]
    resources = ["*"]
  }

  # allow ssm to manage ssh to instances
  statement {
    sid    = "AllowSSMSSHAccess"
    effect = "Allow"
    actions = [
      "ssmmessages:CreateControlChannel",
      "ssmmessages:CreateDataChannel",
      "ssmmessages:OpenControlChannel",
      "ssmmessages:OpenDataChannel"
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "cicd" {
  name        = "QaasRolepolicy"
  description = "A policy granting CodeDeploy permissions"
  policy      = data.aws_iam_policy_document.qaas_cicd_policy.json
  depends_on  = [aws_iam_role.qaas]
}

resource "aws_iam_role_policy_attachment" "attach_custom_policy" {
  depends_on = [aws_iam_policy.cicd]
  role       = aws_iam_role.qaas.name
  policy_arn = aws_iam_policy.cicd.arn
}

resource "aws_iam_instance_profile" "qaas_service" {
  name = "qaas_service"
  path = "/"
  role = aws_iam_role.qaas.name
}

resource "aws_iam_group" "qaas_ssh_group" {
  name = "qaas-ssh-group"
}

resource "aws_iam_group_membership" "qaas_ssh_group" {
  name  = "ssh-group-members"
  users = [var.user_name]
  group = aws_iam_group.qaas_ssh_group.id
}

data "aws_iam_policy_document" "ssh_policy" {
  statement {
    sid       = "AllowUserSSHAccessqaasInstance"
    effect    = "Allow"
    actions   = ["ec2-instance-connect:SendSSHPublicKey"]
    resources = ["arn:aws:ec2:us-west-1:${var.account_id}:instance/*"]
    condition {
      test     = "StringEquals"
      variable = "ec2:osuser"
      values   = ["ec2-user"]
    }
    condition {
      test     = "StringEquals"
      variable = "aws:ResourceTag/qaas"
      values   = ["service"]
    }
  }

  statement {
    sid       = "AllowEC2InstanceConnectCLICalls"
    effect    = "Allow"
    actions   = ["ec2:DescribeInstances"]
    resources = ["*"]
  }
}

resource "aws_iam_group_policy" "ssh_policy" {
  name  = "QaasSSHPolicy"
  group = aws_iam_group.qaas_ssh_group.id

  policy = data.aws_iam_policy_document.ssh_policy.json
}
