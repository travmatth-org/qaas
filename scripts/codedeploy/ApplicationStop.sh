#!/bin/bash
echo "ApplicationStop.sh" | systemd-cat
set -eux pipefail
# send script output to /tmp/ApplicationStop.log for debugging
# exec >> /tmp/ApplicationStop.log 2>&1

if systemctl is-active --quiet httpd; then
  sudo systemctl stop httpd
fi