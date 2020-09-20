data "aws_iam_policy_document" "codedeploy" {
  statement {
    sid     = "QaasCodeDeployTrustRelationships"
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["codedeploy.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "codedeploy_role" {
  name               = "QaasCodeDeployRole"
  assume_role_policy = data.aws_iam_policy_document.codedeploy.json
  description        = "Allows CodeDeploy to call AWS services"
}

resource "aws_iam_role_policy_attachment" "codedeploy_attach" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSCodeDeployRole"
  role       = aws_iam_role.codedeploy_role.name
}