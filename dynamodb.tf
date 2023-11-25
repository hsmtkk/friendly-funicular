resource "aws_dynamodb_table" "table" {
  billing_mode = "PAY_PER_REQUEST"
  name         = "${var.project}-table"
  hash_key     = "link"
  attribute {
    name = "link"
    type = "S"
  }
}