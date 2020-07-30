#!/bin/bash
set -eux pipefail
# send script output to /tmp/ValidateService.log for debugging
# exec >> /tmp/ValidateService.log 2>&1

systemctl is-active --quiet httpd
