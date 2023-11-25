module "rss_read" {
  source        = "./module/lambda"
  filename      = "rss-read"
  function_name = "${var.project}-rss-read"
  role          = aws_iam_role.rss_read.arn
  environment = {
    SQS_QUEUE = aws_sqs_queue.queue.name
  }
}

module "dynamo_write" {
  source        = "./module/lambda"
  filename      = "dynamo-write"
  function_name = "${var.project}-dynamo-write"
  role          = aws_iam_role.dynamo_write.arn
  environment = {
    DYNAMODB_TABLE = aws_dynamodb_table.table.name
  }
}

