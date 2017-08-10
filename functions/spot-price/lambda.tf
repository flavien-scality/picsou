resource "aws_iam_role" "spot-prices-role" {
  name = "spot-prices-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

data "aws_iam_policy_document" "spot-prices-access" {
  statement {
    actions = [
      "ec2:DescribeSpotPriceHistory",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "spot-prices-access" {
  name = "spot-prices-access"
  policy = "${data.aws_iam_policy_document.spot-prices-access.json}"
}

resource "aws_iam_role_policy_attachment" "spot-prices-access" {
  role = "${aws_iam_role.spot-prices-role.name}"
  policy_arn = "${aws_iam_policy.spot-prices-access.arn}"
}

resource "aws_iam_role_policy_attachment" "basic-exec-role" {
  role = "${aws_iam_role.spot-prices-role.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_permission" "event_trigger" {
  statement_id = "AllowExectionFromEvent"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.spot_price.function_name}"
  principal = "events.amazonaws.com"
}

resource "aws_lambda_function" "spot_price" {
  filename = "spot-price.zip"
  function_name = "spot-price"
  handler = "spot-price.handler"
  role = "${aws_iam_role.spot-prices-role.arn}"
  runtime = "python3.6"
  source_code_hash = "${base64sha256(file("spot-price.zip"))}"
  timeout = 10
}
