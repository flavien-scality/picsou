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
