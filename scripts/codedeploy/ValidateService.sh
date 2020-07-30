#!/bin/bash
set -eux pipefail

systemctl is-active --quiet httpd
