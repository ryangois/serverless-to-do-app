resource "aws_dynamodb_table" "to-dos" {
  name = var.dynamo_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "id"
  range_key = "date"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "date"
    type = "S"
  }

  global_secondary_index {
    name = "DateIndex"
    hash_key = "date"
    projection_type = "ALL"
  }
}

