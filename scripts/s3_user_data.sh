#!/bin/bash -xe
# send script output to /tmp/userdata.log for debugging
exec >> /tmp/userdata.log 2>&1

# update yum, install dependencies
sudo yum -y update
sudo yum -y install ruby
sudo yum -y install wget

# install codedeploy agent
cd /home/ec2-user
wget https://aws-codedeploy-us-west-1.s3.amazonaws.com/latest/install
chmod +x ./install
sudo ./intall auto
sudo service codedeploy-agent status
