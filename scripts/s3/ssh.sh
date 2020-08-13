#!/bin/bash

set -eux pipefail
public_ip=`make show.eip`
ssh -i protected/faas_ec2.key ec2-user@${public_ip}