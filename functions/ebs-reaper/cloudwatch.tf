resource "aws_cloudwatch_event_rule" "ebs_reaper_everyday" {
  name = "ebs_reaper_everyday"
  schedule_expression = "rate(3600 minutes)"
}

resource "aws_cloudwatch_event_target" "ebs_reaper_everyday" {
  rule = "${aws_cloudwatch_event_rule.ebs_reaper_everyday.name}"
  target_id = "ebs_reaper"
  arn = "${aws_lambda_function.ebs_reaper.arn}"
}
