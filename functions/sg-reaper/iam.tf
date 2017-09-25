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

data "aws_iam_policy_document" "sg_reaper_access" {
  statement {
    actions = [
      "ec2:DescribeInstances",
      "ec2:DescribeSecurityGroups",
      "ec2:DeleteSecurityGroups",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "sg_reaper_access" {
  name = "sg_reaper_access"
  policy = "${data.aws_iam_policy_document.sg_reaper_access.json}"
}

resource "aws_iam_policy_attachment" "sg_reaper_access" {
  name = "sg_reaper_access"
  roles = ["${aws_iam_role.iam_for_sg_reaper.name}"]
  policy_arn = "${aws_iam_policy.sg_reaper_access.arn}"
}

resource "aws_iam_policy_attachment" "basic_exec_role" {
  name = "basic_exec_role"
  roles = ["${aws_iam_role.iam_for_sg_reaper.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
