resource "aws_cloudwatch_event_rule" "sg_reaper_everyday" {
  name = "sg_reaper_everyday"
  schedule_expression = "rate(3600 minutes)"
}

resource "aws_cloudwatch_event_target" "sg_reaper_everyday" {
  rule = "${aws_cloudwatch_event_rule.sg_reaper_everyday.name}"
  target_id = "sg_reaper"
  arn = "${aws_lambda_function.sg_reaper.arn}"
}
