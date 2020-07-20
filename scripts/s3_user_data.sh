#!/bin/bash

set -ex

# update yum, install dependencies
sudo yum -y update && install ruby wget

# install codedeploy agent
cd /home/ec2-user
wget https://aws-codedeploy-us-west-1.s3.amazonaws.com/latest/install
chmod +x ./install
sudo ./intall auto
sudo service codedeploy-agent status
