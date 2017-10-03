resource "aws_lambda_permission" "ec2_killer" {
  statement_id = "AllowExecutionFromEvent"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ec2_killer.function_name}"
  principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_ec2_killer" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.ec2_killer.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.ec2_killer_everyday.arn}"
}

resource "aws_lambda_function" "ec2_killer" {
  filename = "ec2-killer.zip"
  function_name = "ec2_killer"
  handler = "ec2-killer.handler"
  role = "${aws_iam_role.iam_for_ec2_killer.arn}"
  runtime = "python3.6"
  source_code_hash = "${base64sha256(file("ec2-killer.zip"))}"
  timeout = 60
}
