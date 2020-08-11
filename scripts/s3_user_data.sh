#!/bin/bash
# shellcheck disable=SC2154
set -eux pipefail
# send script output to /tmp/user_data.log for debugging
exec >> /tmp/user_data.log 2>&1

# install dependencies
sudo yum -y update
sudo yum -y install ruby wget libcap2-bin shadow-utils.x86_64

# install codedeploy agent
cd /home/ec2-user
wget https://aws-codedeploy-us-west-1.s3.us-west-1.amazonaws.com/latest/install
chmod +x ./install
sudo ./install auto
sudo service codedeploy-agent status

# create user for faas service
sudo useradd -s /sbin/nologin -M faas
# create dirs for static content
sudo mkdir -p web/www/static
# create log dir
sudo mkdir /var/log/httpd

# install cloudwatch-agent
sudo yum install -y epel-release
sudo yum --enablerepo=epel install collectd
cw_config="/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json"
mkdir -P /tmp/cloudwatch-logs
cd /tmp/cloudwatch-logs
wget https://s3.us-west-1.amazonaws.com/amazoncloudwatch-agent-us-west-1/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm
sudo rpm -U ./amazon-cloudwatch-agent.rpm
# https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-Agent-Configuration-File-Details.html
cat <<EOF > $cw_config
{
   "agent": {
      "metrics_collection_interval": 15
   },
   "metrics": {
      "namespace": "Prod/Faas",
      "metrics_collected": {
         "cpu": {
            "measurement": [
               "usage_active",
               "usage_iowait"
            ],
         },
         "mem": {
            "measurement": [
               "used",
               "total"
            ],
         },
         "net": {
            "measurement": [
               "bytes_sent",
               "bytes_recv"
            ],
         },
      },
      "append_dimensions": {
         "ImageId": "${aws:ImageId}",
         "InstanceId": "${aws:InstanceId}",
         "InstanceType": "${aws:InstanceType}",
         "AutoScalingGroupName": "${aws:AutoScalingGroupName}"
      },
      "aggregation_dimensions": [
         [
            "AutoScalingGroupName"
         ],
         [
            "InstanceId",
            "InstanceType"
         ]
      ]
   },
   "logs": {
      "logs_collected": {
         "files": {
            "collect_list": [
               {
                  "file_path": "/var/log/httpd.log*",
                  "log_group_name": "faas-httpd-logs",
                  "log_stream_name": "ec2-${instance_id}-logs",
               }
            ]
         }
      },
      "log_stream_name": "${instance_id}/log-stream"
   }
}
EOF

# start the cloudwatch agent
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
   -a fetch-config \
   -m ec2 \
   -c file:$cw_config \
   -s

# veryify the cloudwatch agent is running
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
   -m ec2 \
   -a status
