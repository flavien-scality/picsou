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
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "iam_policy_for_ec2_killer" {
  name = "iam_policy_for_ec2_killer"
  role = "${aws_iam_role.iam_for_ec2_killer.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "ec2:TerminateInstances",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
