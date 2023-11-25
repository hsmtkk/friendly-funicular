resource "aws_scheduler_schedule_group" "group" {
  name = "${var.project}-group"
}

resource "aws_scheduler_schedule" "schedule" {
  name       = "${var.project}-schedule"
  group_name = aws_scheduler_schedule_group.group.name

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression = "rate(1 hours)"

  target {
    arn      = module.rss_read.arn
    role_arn = aws_iam_role.eventbridge.arn
  }
}
