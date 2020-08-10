resource "aws_cloudwatch_log_group" "faas" {
	name = "faas_log_group"

	tags = {
		faas = "service"
	}
}

resource "aws_cloudwatch_dashboard" "faas-dashboard" {
	dashboard_name = "dashboard-${aws_instance.faas_service.id}"

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
					"title": "CPU Utilization"
				}
			}
		]
	}
	EOF
}
