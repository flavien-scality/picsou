resource "aws_iam_role" "spot_prices_role" {
  name = "spot_prices_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

data "aws_iam_policy_document" "spot_prices_access_doc" {
  statement {
    actions = [
      "ec2:DescribeSpotPriceHistory",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "spot_prices_access_policy" {
  name = "spot_prices_access_policy"
  policy = "${data.aws_iam_policy_document.spot_prices_access_doc.json}"
}

resource "aws_iam_role_policy_attachment" "spot_prices_access_attach" {
  role = "${aws_iam_role.spot_prices_role.name}"
  policy_arn = "${aws_iam_policy.spot_prices_access_policy.arn}"
}

resource "aws_iam_role_policy_attachment" "basic_exec_role" {
  role = "${aws_iam_role.spot_prices_role.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
