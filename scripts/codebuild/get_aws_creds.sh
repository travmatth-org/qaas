#!/bin/bash
# "> 'set -o pipefail on' is the same as 'set -o pipefail; set on', which"
# "> turns it on but also changes $*."
# http://gnu-bash.2382.n7.nabble.com/set-command-overrides-my-ARGV-array-td18459.html
role_arn="$1"
set -eux pipefail

role_session_name='packer'

temp_role=$(aws sts assume-role \
     --role-arn $role_arn \
     --role-session-name $role_session_name \
     --output json)

export AWS_ACCESS_KEY_ID=$(echo $temp_role | jq -r .Credentials.AccessKeyId)
export AWS_SECRET_ACCESS_KEY=$(echo $temp_role | jq -r .Credentials.SecretAccessKey)
export AWS_SESSION_TOKEN=$(echo $temp_role | jq -r .Credentials.SessionToken)
export AWS_DEFAULT_REGION=us-west-1

aws configure set aws_access_key_id $AWS_ACCESS_KEY_ID
aws configure set aws_secret_access_key $AWS_SECRET_ACCESS_KEY
aws configure set aws_session_token $AWS_SESSION_TOKEN
aws configure set region $AWS_DEFAULT_REGION