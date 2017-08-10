resource "aws_volume_attachment" "ebs_setup" {
  device_name = "/dev/sdh"
  volume_id = "${aws_ebs_volume.setup-dev.id}"
  instance_id = "${aws_instance.setup-dev.id}"
}

resource "aws_ebs_volume" "setup-dev" {
  availability_zone = "eu-west-2a"
  size = 40
  type = "gp2"
  tags = {
    Name = "docker-build-farmer"
    Team = "setup"
    Spawner = "terraform"
  }
}
