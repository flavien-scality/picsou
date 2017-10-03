resource "aws_lambda_permission" "sg_reaper" {
  statement_id = "AllowExecutionFromEvent"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.sg_reaper.function_name}"
  principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_sg_reaper" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.sg_reaper.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.sg_reaper_everyday.arn}"
}

resource "aws_lambda_function" "sg_reaper" {
  filename = "sg-reaper.zip"
  function_name = "sg_reaper"
  handler = "sg-reaper.handler"
  role = "${aws_iam_role.iam_for_sg_reaper.arn}"
  runtime = "python3.6"
  source_code_hash = "${base64sha256(file("sg-reaper.zip"))}"
  timeout = 60
}
