resource "aws_s3_bucket" "bootstrap_scripts" {
  bucket = "${var.login}-dev"
  acl = "private"
}

resource "aws_s3_bucket_object" "dev_bootstrap_script" {
  bucket = "${var.login}-dev"
  key = "bootstrap-scripts/dev_bootstrap_script.sh"
  source = "shared/dev_script.sh"
}
