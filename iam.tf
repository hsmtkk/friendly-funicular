data "aws_iam_policy_document" "eventbridge" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["scheduler.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "eventbridge" {
  name                = "${var.project}-eventbridge"
  assume_role_policy  = data.aws_iam_policy_document.eventbridge.json
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaRole"]
}

data "aws_iam_policy_document" "lambda" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "dynamo_write" {
  name                = "${var.project}-dynamo-write"
  assume_role_policy  = data.aws_iam_policy_document.lambda.json
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole", "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole", "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"]
}

resource "aws_iam_role" "rss_read" {
  name                = "${var.project}-rss-read"
  assume_role_policy  = data.aws_iam_policy_document.lambda.json
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole", "arn:aws:iam::aws:policy/AmazonSQSFullAccess"]
}
