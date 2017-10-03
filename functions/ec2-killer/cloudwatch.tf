resource "aws_cloudwatch_event_rule" "ec2_killer_everyday" {
  name = "ec2_killer_everyday"
  schedule_expression = "rate(3600 minutes)"
}

resource "aws_cloudwatch_event_target" "ec2_killer_everyday" {
  rule = "${aws_cloudwatch_event_rule.ec2_killer_everyday.name}"
  target_id = "ec2-killer"
  arn = "${aws_lambda_function.ec2_killer.arn}"
}
