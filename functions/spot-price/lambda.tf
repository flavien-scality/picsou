resource "aws_lambda_permission" "event_trigger" {
  statement_id = "AllowExectionFromEvent"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.spot_price.function_name}"
  principal = "events.amazonaws.com"
}

resource "aws_lambda_function" "spot_price" {
  filename = "spot-price.zip"
  function_name = "spot_price"
  handler = "spot-price.handler"
  role = "${aws_iam_role.spot_prices_role.arn}"
  runtime = "python3.6"
  source_code_hash = "${base64sha256(file("spot-price.zip"))}"
  timeout = 10
}
