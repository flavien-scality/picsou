resource "aws_lambda_function" "spot_price" {
  filename         = "spot-price.zip"
  function_name    = "spot-price"
  role             = "${aws_iam_role.iam_for_spot_price.arn}"
  handler          = "main.handle"
  source_code_hash = "${base64sha256(file("spot-price.zip"))}"
  runtime          = "python2.7"
}
