resource "aws_sqs_queue" "queue" {
  name = "${var.project}-queue"
}