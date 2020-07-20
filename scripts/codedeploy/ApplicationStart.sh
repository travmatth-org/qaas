#!/bin/bash

set -ex

sudo systemctl start httpd
sudo systemctl enable httpd