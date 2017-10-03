resource "aws_lambda_permission" "ebs_reaper" {
  statement_id = "AllowExecutionFromEvent"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ebs_reaper.function_name}"
  principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_ebs_reaper" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.ebs_reaper.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.ebs_reaper_everyday.arn}"
}

resource "aws_lambda_function" "ebs_reaper" {
  filename = "ebs-reaper.zip"
  function_name = "ebs_reaper"
  handler = "ebs-reaper.handler"
  role = "${aws_iam_role.iam_for_ebs_reaper.arn}"
  runtime = "python3.6"
  source_code_hash = "${base64sha256(file("ebs-reaper.zip"))}"
  timeout = 60
}
