#!/bin/bash -xe

# send script output to /tmp/userdata.log for debugging
exec >> /tmp/userdata.log 2>&1

# install dependencies
sudo yum -y update
sudo yum -y install ruby wget libcap2-bin

# install codedeploy agent
cd /home/ec2-user
wget https://aws-codedeploy-us-west-1.s3.amazonaws.com/latest/install
chmod +x ./install
sudo ./intall auto
sudo service codedeploy-agent status

# create user for faas service
sudo useradd faas -s /sbin/nogin -M

# allow faas service to run on privileged port
sudo setcap 'cap_net_bind_service=+ep' /usr/sbin/httpd
