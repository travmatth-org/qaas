resource "aws_dynamodb_table" "qaas_quote_table" {
  name           = "qaas-quote-table"
  billing_mode   = "PAY_PER_REQUEST"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }

  tags = {
    qaas = "dynamodb"
  }
}

resource "aws_dynamodb_table" "qaas_author_table" {
  name           = "qaas-author-table"
  billing_mode   = "PAY_PER_REQUEST"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Name"

  attribute {
    name = "Name"
    type = "S"
  }

  tags = {
    qaas = "dynamodb"
  }
}

resource "aws_dynamodb_table" "qaas_topics_table" {
  name           = "qaas-topics-table"
  billing_mode   = "PAY_PER_REQUEST"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Topic"

  attribute {
    name = "Topic"
    type = "S"
  }

  tags = {
    qaas = "dynamodb"
  }
}