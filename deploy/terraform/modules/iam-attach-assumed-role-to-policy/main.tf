resource "aws_iam_role_policy_attachment" "tf_iam_attach_assumed_role_to_permissions_policy" {
  role       = aws_iam_role.tf_iam_assumed_role.name
  policy_arn = aws_iam_policy.tf_iam_assumed_policy.arn

  lifecycle {
    prevent_destroy = true
  }
}