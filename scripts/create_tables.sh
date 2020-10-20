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
	--table-name qaas-quote-table \
	--attribute-definitions \
		AttributeName=ID,AttributeType=S \
	--key-schema \
		AttributeName=ID,KeyType=HASH \
	--provisioned-throughput \
		ReadCapacityUnits=1,WriteCapacityUnits=1

echo "Creating qaas author table"
aws dynamodb --endpoint-url "$endpoint" create-table \
	--table-name qaas-author-table \
	--attribute-definitions \
		AttributeName=Name,AttributeType=S \
		AttributeName=QuoteID,AttributeType=S \
	--key-schema \
		AttributeName=Name,KeyType=HASH \
		AttributeName=QuoteID,KeyType=RANGE \
	--provisioned-throughput \
		ReadCapacityUnits=1,WriteCapacityUnits=1

echo "Creating qaas topics table"
aws dynamodb --endpoint-url "$endpoint" create-table \
	--table-name qaas-topic-table \
	--attribute-definitions \
		AttributeName=Name,AttributeType=S \
		AttributeName=QuoteID,AttributeType=S \
	--key-schema \
		AttributeName=Name,KeyType=HASH \
		AttributeName=QuoteID,KeyType=RANGE \
	--provisioned-throughput \
		ReadCapacityUnits=1,WriteCapacityUnits=1
