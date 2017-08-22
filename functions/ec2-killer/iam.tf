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

resource "aws_lambda_permission" "ec2_killer" {
  statement_id = "AllowExecutionFromEvent"
  action = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ec2_killer.function_name}"
  principal = "events.amazonaws.com"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_ec2_killer" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.ec2_killer.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.everyday.arn}"
}

resource "aws_lambda_function" "ec2_killer" {
  filename = "ec2-killer.zip"
  function_name = "ec2_killer"
  handler = "ec2-killer.handler"
  role = "${aws_iam_role.iam_for_ec2_killer.arn}"
  runtime = "python3.6"
  source_code_hash = "${base64sha256(file("ec2-killer.zip"))}"
  timeout = 60
}
