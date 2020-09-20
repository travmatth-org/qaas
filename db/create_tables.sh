#!/bin/bash
# "> 'set -o pipefail on' is the same as 'set -o pipefail; set on', which"
# "> turns it on but also changes $*."
# http://gnu-bash.2382.n7.nabble.com/set-command-overrides-my-ARGV-array-td18459.html
endpoint="$1"
set -eux pipefail

# https://github.com/aws/aws-cli/pull/4702#issue-344978525
export AWS_PAGER="" 

echo "Creating qaas quotes table"
aws dynamodb --endpoint-url "$endpoint" create-table \
	--table-name qaas-quotes \
	--attribute-definitions "AttributeName=Id,AttributeType=S" \
	--key-schema "AttributeName=Id,KeyType=HASH" \
	--provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"

echo "Creating qaas author table"
aws dynamodb --endpoint-url "$endpoint" create-table \
	--table-name qaas-author \
	--attribute-definitions "AttributeName=Name,AttributeType=S" \
	--key-schema "AttributeName=Name,KeyType=HASH" \
	--provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"

echo "Creating qaas topics table"
aws dynamodb --endpoint-url "$endpoint" create-table \
	--table-name qaas-topics \
	--attribute-definitions "AttributeName=Topic,AttributeType=S" \
	--key-schema "AttributeName=Topic,KeyType=HASH" \
	--provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"
