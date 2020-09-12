resource "aws_cloudwatch_log_group" "faas" {
	name = "faas-httpd-logs"

	tags = {
		faas = "service"
	}
}

resource "aws_cloudwatch_log_stream" "foo" {
  name           = "${aws_autoscaling_group.faas_service.name}-logs"
  log_group_name = aws_cloudwatch_log_group.faas.name
}

resource "aws_cloudwatch_dashboard" "faas-dashboard" {
	dashboard_name = "dashboard-faas-service-${aws_autoscaling_group.faas_service.name}"

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
						"${aws_autoscaling_group.faas_service.name}"
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
							"${aws_autoscaling_group.faas_service.name}"
						],
						[
							"AWS/EC2",
							"NetworkOut",
							"AutoScalingGroupName",
							"${aws_autoscaling_group.faas_service.name}"
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
# 	alarm_name			= "faas-alarm-high-cpu"
# 	comparison_operator	= "GreaterThanOrEqualToThreshold"
# 	evaluation_periods	= "2"
# 	metric_name			= "CPUUtilization"
# 	namespace			= "AWS/EC2"
# 	period				= "120"
# 	statistic			= "Average"
# 	threshold			= "70"

# 	dimensions			= {
# 		AutoScalingGroupName = aws_autoscaling_group.faas_service.name
# 	}

# 	alarm_description	= "Monitor EC2 instance CPU utilization, shutdown if average >= 70%"
# 	alarm_actions		= [aws_autoscaling_group.faas_service.name]
# }