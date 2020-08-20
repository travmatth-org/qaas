#!/bin/bash
set -eux pipefail
# send script output to /tmp/user_data.log for debugging
exec >> /tmp/user_data.log 2>&1

# install dependencies
sudo yum -y update
sudo yum -y install ruby wget shadow-utils.x86_64

# install codedeploy agent
cd /home/ec2-user
wget https://aws-codedeploy-us-west-1.s3.us-west-1.amazonaws.com/latest/install
chmod +x ./install
sudo ./install auto
sudo service codedeploy-agent status
sudo rm install

# create user for faas service
sudo useradd -s /sbin/nologin -M faas

# create dirs for static content
sudo mkdir -p web/www/static
# create log dir
sudo mkdir /var/log/httpd

# install cloudwatch-agent
mkdir -p /tmp/cloudwatch-logs
cd /tmp/cloudwatch-logs
wget https://s3.us-west-1.amazonaws.com/amazoncloudwatch-agent-us-west-1/amazon_linux/amd64/latest/amazon-cloudwatch-agent.rpm
sudo rpm -U ./amazon-cloudwatch-agent.rpm

# add config file for cw agent, specifying metrics & logs to collect
# https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-Agent-Configuration-File-Details.html
cat << 'EOF' | sudo tee /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json
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
            ]
         },
         "mem": {
            "measurement": [
               "used",
               "total"
            ]
         },
         "net": {
            "measurement": [
               "bytes_sent",
               "bytes_recv"
            ]
         }
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
   }
}
EOF

# start the cloudwatch agent
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
   -a fetch-config \
   -m ec2 \
   -c file:/opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json \
   -s

# verify the cloudwatch agent is running
sudo /opt/aws/amazon-cloudwatch-agent/bin/amazon-cloudwatch-agent-ctl \
   -m ec2 \
   -a status

# install x-ray daemon
curl https://s3.dualstack.us-east-2.amazonaws.com/aws-xray-assets.us-east-2/xray-daemon/aws-xray-daemon-3.x.rpm -o /home/ec2-user/xray.rpm
sudo yum install -y /home/ec2-user/xray.rpm

# start x-ray daemon
sudo systemctl enable xray
sudo systemctl start xray

# verify the xray agent is running
sudo systemctl is-active --quiet xray

# install vector to forward journald traffic to the cloudwatch agent
mkdir -p /tmp/vector
curl -O https://packages.timber.io/vector/0.10.X/vector-x86_64.rpm
sudo yum install -y vector-x86_64.rpm

# configure vector daemon  
cat << VECTORCFG | sudo tee /etc/vector/vector.toml
[sources.in]
  type = "journald"
  include_units = ["httpd"]

[sinks.out]
  encoding.codec = "json"

  # General
  group_name = "faas-httpd-logs"
  inputs = ["in"]
  region = "us-west-1"
  stream_name = "{{ host }}"
  type = "aws_cloudwatch_logs"
VECTORCFG

# start vector daemon
sudo usermod -aG systemd-journal vector
sudo systemctl enable vector
sudo systemctl start vector

# verify the vector daemon is running
sudo systemctl is-active --quiet vector

# add functions to bashrc
cat << BASHRC > /home/ec2-user/.bashrc
# https://stackoverflow.com/questions/23151425/how-to-run-cloud-init-manually/23152036
redo-cloud-init() {
	sudo cloud-init clean;
	sudo cloud-init start;
}

alias ..="cd .."
alias=logs="journalctl -xe --no-pager"
alias httpd-status="sudo systemctl status httpd"
alias httpd-logs="journalctl -xe --no-pager -u httpd"
alias fhere="find . -name"
BASHRC