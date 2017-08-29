resource "aws_cloudwatch_event_rule" "everyday" {
  name = "everyday"
  schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "sg_reaper_everyday" {
  rule = "${aws_cloudwatch_event_rule.everyday.name}"
  target_id = "sg-reaper"
  arn = "${aws_lambda_function.sg_reaper.arn}"
}
