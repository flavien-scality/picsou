provider "aws" {
  region = "eu-west-2"
}

data "aws_ami" "testing" {
  most_recent = true

  filter {
    name   = "name"
    values = ["testing-*"]
  }

  owners = ["self"]
}

resource "aws_instance" "testing" {
  ami           = "${data.aws_ami.testing.id}"
  instance_type = "t2.micro"
  iam_instance_profile = "${aws_iam_instance_profile.testing_instance_profile.name}"
  key_name = "maxime-london"
  tags {
    Name = "testing-maxime"
  }
}
