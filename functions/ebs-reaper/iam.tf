resource "aws_iam_role" "iam_for_ebs_reaper" {
  name = "iam_for_ebs_reaper"

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

data "aws_iam_policy_document" "ebs_reaper_access_document" {
  statement {
    actions = [
      "ec2:DeleteVolume",
      "ec2:DescribeVolumes",
    ]
    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "ebs_reaper_access_policy" {
  name = "ebs_reaper_access_policy"
  policy = "${data.aws_iam_policy_document.ebs_reaper_access_document.json}"
}

resource "aws_iam_policy_attachment" "ebs_reaper_access_attach" {
  name = "ebs_reaper_access_attach"
  roles = ["${aws_iam_role.iam_for_ebs_reaper.name}"]
  policy_arn = "${aws_iam_policy.ebs_reaper_access_policy.arn}"
}

resource "aws_iam_policy_attachment" "basic_exec_role" {
  name = "basic_exec_role"
  roles = ["${aws_iam_role.iam_for_ebs_reaper.name}"]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
