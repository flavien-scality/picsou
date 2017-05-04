resource "aws_lambda_function" "ec2_killer" {
  filename         = "ec2-killer.zip"
  function_name    = "ec2-killer"
  role             = "${aws_iam_role.iam_for_ec2_killer.arn}"
  handler          = "main.handle"
  source_code_hash = "${base64sha256(file("ec2-killer.zip"))}"
  runtime          = "python2.7"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_ec2_killer" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ec2_killer.function_name}"
  principal = "events.amazonaws.com"
  source_arn = "${aws_cloudwatch_event_rule.everyday.arn}"
}
