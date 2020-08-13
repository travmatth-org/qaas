resource "aws_cloudwatch_log_group" "faas" {
	name = "faas-httpd-logs"

	tags = {
		faas = "service"
	}
}

resource "aws_cloudwatch_log_stream" "foo" {
  name           = "ec2-${aws_instance.faas_service.id}-logs"
  log_group_name = aws_cloudwatch_log_group.faas.name
}

resource "aws_cloudwatch_dashboard" "faas-dashboard" {
	dashboard_name = "dashboard-faas-service-ec2-${aws_instance.faas_service.id}"

	# https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/CloudWatch-Dashboard-Body-Structure.html#CloudWatch-Dashboard-Properties-Metrics-Array-Format
	dashboard_body = <<-EOF
	{
		"widgets": [
			{
				"type": "metric",
				"properties": {
					"metrics": [[
						"AWS/EC2",
						"CPUUtilization",
						"InstanceId",
						"${aws_instance.faas_service.id}"
					]],
					"period": 300,
					"stat": "Average",
					"region": "us-west-1",
					"title": "CPU Utilization"
				}
			},
			{
				"type": "metric",
				"properties": {
					"metrics": [
						[
							"AWS/EC2",
							"NetworkIn",
							"InstanceId",
							"${aws_instance.faas_service.id}"
						],
						[
							"AWS/EC2",
							"NetworkOut",
							"InstanceId",
							"${aws_instance.faas_service.id}"
						]
					],
					"period": 300,
					"stat": "Average",
					"region": "us-west-1",
					"title": "Network Traffic"
				}
			}
		]
	}
	EOF
}
