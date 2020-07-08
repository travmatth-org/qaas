resource "aws_iam_role_policy" "codebuild_iam_role_policy" {
  name = var.codebuild_iam_role_policy_name
  role = var.codebuild_iam_role_name

  policy = <<POLICY
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
        "${var.logging_bucket_arn}",
        "${var.logging_bucket_arn}/*",
        "${var.state_bucket_arn}",
        "${var.state_bucket_arn}/*",
        "arn:aws:s3:::codepipeline-us-east-1*",
        "arn:aws:s3:::codepipeline-us-east-1*/*",
        "${var.tf_codepipeline_artifact_bucket_arn}",
        "${var.tf_codepipeline_artifact_bucket_arn}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:*"
      ],

      "Resource": "${var.dynamodb_lock_state_table_arn}"
    },
    {
      "Effect": "Allow",
      "Action": [
        "iam:Get*",
        "iam:List*"
      ],
      "Resource": "${var.codebuild_iam_role_arn}"
    },
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Resource": "${var.codebuild_iam_role_arn}"
    }
  ]
}
POLICY
}
