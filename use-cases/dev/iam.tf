resource "aws_iam_instance_profile" "dev-profile" {
  name = "dev-profile"
  roles = [ "${aws_iam_role.dev-role.name}" ]
}

resource "aws_iam_role" "dev-role" {
  name = "dev-role"
  assume_role_policy = <<EOF
{
  "Version": "2008-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "s3-policy" {
  name = "s3-policy"
  role = "${aws_iam_role.dev-role.id}"
  policy = <<EOF
{
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "s3:*"
    ],
    "Resource": "arn:aws:s3:::mvaude-dev/*"
  }]
}
EOF
}
