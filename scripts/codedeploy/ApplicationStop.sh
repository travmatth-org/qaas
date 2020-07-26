#!/bin/bash

set -eux pipefail

# send script output to /tmp/ApplicationStop.log for debugging
exec >> /tmp/ApplicationStop.log 2>&1

systemctl is-active --quiet httpd

if [ $? -neq 0]; then
  sudo systemctl stop httpd
fi