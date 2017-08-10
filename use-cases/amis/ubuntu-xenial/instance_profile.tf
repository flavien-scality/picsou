resource "aws_iam_instance_profile" "testing_instance_profile" {
  name  = "testing-instance-profile"
  roles = ["${aws_iam_role.testing_role.name}"]
}

resource "aws_iam_role" "testing_role" {
  name = "testing-role"

  assume_role_policy = <<EOF
{
  "Version": "2017-04-12",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "CloudWatchAccess" {
  name        = "CloudWatchAccess-testing"

  policy = <<EOF
{
    "Version": "2017-04-12",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
               "cloudwatch:EnableAlarmActions",
               "cloudwatch:GetMetricData",
               "cloudwatch:GetMetricStatistics",
               "cloudwatch:ListMetrics",
               "cloudwatch:PutMetricAlarm",
               "cloudwatch:PutMetricData",
               "cloudwatch:SetAlarmState",
               "logs:CreateLogGroup",
               "logs:CreateLogStream",
               "logs:GetLogEvents",
               "logs:PutLogEvents",
               "logs:DescribeLogGroups",
               "logs:DescribeLogStreams",
               "logs:PutRetentionPolicy"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_policy_attachment" "attach_cloudwatch" {
  name       = "testing-iam-attachment"
  policy_arn = "${aws_iam_policy.CloudWatchAccess.arn}"
  roles      = ["${aws_iam_role.testing_role.name}"]
}
