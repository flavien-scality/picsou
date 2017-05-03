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

resource "aws_lambda_function" "ec2_killer" {
  filename         = "ec2-killer.zip"
  function_name    = "ec2-killer"
  role             = "${aws_iam_role.iam_for_ec2_killer.arn}"
  handler          = "main.handle"
  source_code_hash = "${base64sha256(file("ec2-killer.zip"))}"
  runtime          = "python2.7"
}

resource "aws_cloudwatch_event_rule" "everyday" {
  name = "everyday"
  schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "ec2_killer_everyday" {
  rule = "${aws_cloudwatch_event_rule.everyday.name}"
  target_id = "ec2-killer"
  arn = "${aws_lambda_function.ec2_killer.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_ec2_killer" {
  statement_id = "AllowExecutionFromCloudWatch"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ec2_killer.function_name}"
  principal = "events.amazonaws.com"
  source_arn = "${aws_cloudwatch_event_rule.everyday.arn}"
}
