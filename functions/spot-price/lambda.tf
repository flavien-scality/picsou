resource "aws_lambda_function" "spot_price" {
  filename         = "spot-price.zip"
  function_name    = "spot-price"
  handler          = "main.handle"
  source_code_hash = "${base64sha256(file("spot-price.zip"))}"
  runtime          = "python2.7"
}
