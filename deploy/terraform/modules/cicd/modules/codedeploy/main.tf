resource "aws_codedeploy_app" "qaas" {
  name = "qaas"
}

resource "aws_codedeploy_deployment_group" "qaas_in_place" {
  app_name              = aws_codedeploy_app.qaas.name
  deployment_group_name = "${aws_codedeploy_app.qaas.name}-deployment-group"
  service_role_arn      = aws_iam_role.codedeploy_role.arn

  ec2_tag_set {
    ec2_tag_filter {
      type  = "KEY_AND_VALUE"
      key   = "qaas"
      value = "service"
    }
  }

  deployment_style {
    deployment_type   = "IN_PLACE"
    deployment_option = "WITHOUT_TRAFFIC_CONTROL"
  }
}

output "app" {
  value = aws_codedeploy_app.qaas
}

output "deployment_group" {
  value = aws_codedeploy_deployment_group.qaas_in_place
}
