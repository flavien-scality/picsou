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

data "aws_iam_policy_document" "ec2_killer_access" {
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

resource "aws_iam_policy" "ec2_killer_access" {
  name = "ec2_killer_access"
  policy = "${data.aws_iam_policy_document.ec2_killer_access.json}"
}

resource "aws_iam_policy_attachment" "ec2_killer_access" {
  name = "ec2_killer_access"
  roles = ["${aws_iam_role.iam_for_ec2_killer.name}"]
  policy_arn = "${aws_iam_policy.ec2_killer_access.arn}"
}

resource "aws_iam_policy_attachment" "basic_exec_role" {
  name = "basic_exec_role"
  roles = ["${aws_iam_role.iam_for_ec2_killer.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
