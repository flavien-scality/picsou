resource "aws_iam_role" "iam_for_sg_reaper" {
  name = "iam_for_sg_reaper"

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

data "aws_iam_policy_document" "sg-reaper-access" {
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

resource "aws_iam_policy" "sg-reaper-access" {
  name = "sg-reaper-access"
  policy = "${data.aws_iam_policy_document.sg-reaper-access.json}"
}

resource "aws_iam_policy_attachment" "sg-reaper-access" {
  name = "sg-reaper-access"
  roles = ["${aws_iam_role.iam_for_sg_reaper.name}"]
  policy_arn = "${aws_iam_policy.sg-reaper-access.arn}"
}

resource "aws_iam_policy_attachment" "basic-exec-role" {
  name = "basic-exec-role"
  roles = ["${aws_iam_role.iam_for_sg_reaper.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
