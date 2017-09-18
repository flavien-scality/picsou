resource "aws_iam_role" "iam_for_ec2_killer" {
  name = "iam_for_ec2_killer"

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

data "aws_iam_policy_document" "ec2-killer-access" {
  statement {
    actions = [
      "ec2:TerminateInstances",
      "ec2:DescribeInstances",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "ec2-killer-access" {
  name = "ec2-killer-access"
  policy = "${data.aws_iam_policy_document.ec2-killer-access.json}"
}

resource "aws_iam_policy_attachment" "ec2-killer-access" {
  name = "ec2-killer-access"
  roles = ["${aws_iam_role.iam_for_ec2_killer.name}"]
  policy_arn = "${aws_iam_policy.ec2-killer-access.arn}"
}

resource "aws_iam_policy_attachment" "basic-exec-role" {
  name = "basic-exec-role"
  roles = ["${aws_iam_role.iam_for_ec2_killer.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
