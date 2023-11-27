module "rss_read" {
  source        = "./module/lambda"
  filename      = "rss-read"
  function_name = "${var.project}-rss-read"
  role          = aws_iam_role.rss_read.arn
  environment = {
    NEWS_URL  = "https://news.yahoo.co.jp/rss/topics/top-picks.xml"
    QUEUE_URL = aws_sqs_queue.queue.name
    QUEUE_URL = "https://sqs.${var.region}.amazonaws.com/${data.aws_caller_identity.current.account_id}/${aws_sqs_queue.queue.name}"
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

