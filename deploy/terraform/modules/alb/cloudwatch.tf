resource "aws_cloudwatch_log_group" "qaas" {
  name = "qaas-httpd-logs"

  tags = {
    qaas = "service"
  }
}

resource "aws_cloudwatch_log_stream" "foo" {
  name           = "${aws_autoscaling_group.qaas_service.name}-logs"
  log_group_name = aws_cloudwatch_log_group.qaas.name
}

resource "aws_cloudwatch_dashboard" "qaas-dashboard" {
  dashboard_name = "dashboard-qaas-service-${aws_autoscaling_group.qaas_service.name}"

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
						"AutoScalingGroupName",
						"${aws_autoscaling_group.qaas_service.name}"
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
							"AutoScalingGroupName",
							"${aws_autoscaling_group.qaas_service.name}"
						],
						[
							"AWS/EC2",
							"NetworkOut",
							"AutoScalingGroupName",
							"${aws_autoscaling_group.qaas_service.name}"
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

# resource "aws_cloudwatch_metric_alarm" "high_cpu" {
# 	alarm_name			= "qaas-alarm-high-cpu"
# 	comparison_operator	= "GreaterThanOrEqualToThreshold"
# 	evaluation_periods	= "2"
# 	metric_name			= "CPUUtilization"
# 	namespace			= "AWS/EC2"
# 	period				= "120"
# 	statistic			= "Average"
# 	threshold			= "70"

# 	dimensions			= {
# 		AutoScalingGroupName = aws_autoscaling_group.qaas_service.name
# 	}

# 	alarm_description	= "Monitor EC2 instance CPU utilization, shutdown if average >= 70%"
# 	alarm_actions		= [aws_autoscaling_group.qaas_service.name]
# }