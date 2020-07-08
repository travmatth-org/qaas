resource "aws_iam_policy" "tf_iam_assumed_policy" {
  name = "TerraformAssumedIamPolicy"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowAllPermissions",
      "Effect": "Allow",
      "Action": [
        "*"
      ],
      "Resource": "*"
    }
  ]
}
EOF

  lifecycle {
    prevent_destroy = true
  }
}
