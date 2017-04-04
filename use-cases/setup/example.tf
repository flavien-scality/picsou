resource "aws_security_group" "setup-group" {
  name = "setup-sg"
  description = "Setup Security Group"
}

resource "aws_security_group_rule" "ssh_ingress_access" {
  type = "ingress"
  from_port = 22
  to_port = 22
  protocol = "tcp"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.setup-group.id}"
}

resource "aws_security_group_rule" "egress_access" {
  type = "egress"
  from_port = 0
  to_port = 0
  protocol = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.setup-group.id}"
}

resource "aws_instance" "setup-dev" {
  ami = "ami-ede2e889"
  root_block_device {
    volume_size = "100"
    volume_type = "gp2"
    delete_on_termination = "true"
  }
  ebs_block_device {
    device_name = "/dev/sdb"
    volume_size = "50"
    volume_type = "gp2"
    delete_on_termination = "true"
  }
  instance_type = "t2.2xlarge"
  vpc_security_group_ids = [ "${aws_security_group.setup-group.id}" ]
  tags {
    Name = "docker-build-farmer"
    Team = "setup"
  }
  user_data = "${file("shared/user-data.txt")}"
}
