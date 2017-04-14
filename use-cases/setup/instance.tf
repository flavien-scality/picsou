resource "aws_instance" "setup-dev" {
  ami = "${var.aws_amis[var.os]}"
  availability_zone = "eu-west-2a"
  root_block_device {
    volume_size = "100"
    volume_type = "gp2"
    delete_on_termination = "true"
  }
  instance_type = "t2.2xlarge"
  vpc_security_group_ids = [ "${aws_security_group.setup-group.id}" ]
  tags {
    Name = "docker-build-farmer"
    Team = "setup"
    Spawner = "terraform"
  }
  user_data = "${file("shared/user-data.txt")}"
}
