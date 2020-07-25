#!/bin/bash

set -eux pipefail

# send script output to /tmp/ApplicationStart.log for debugging
exec >> /tmp/ApplicationStart.log 2>&1

sudo systemctl start httpd
