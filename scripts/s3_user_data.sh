#!/bin/bash
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
