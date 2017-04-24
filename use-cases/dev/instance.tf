resource "aws_instance" "dev-machine" {
  ami = "${var.aws_amis[var.os]}"
  iam_instance_profile = "${aws_iam_instance_profile.dev-profile.name}"
  availability_zone = "eu-west-2a"
  monitoring = true
  depends_on = ["aws_s3_bucket_object.dev_bootstrap_script"]
  root_block_device {
    volume_size = "100"
    volume_type = "gp2"
    delete_on_termination = "true"
  }
  instance_type = "t2.nano"
  vpc_security_group_ids = [ "${aws_security_group.dev-group.id}" ]
  tags {
    Name = "dev-machine"
    User = "${var.login}"
    Spawner = "terraform"
  }
  user_data = "${file("shared/user-data.txt.tmpl")}"
}
