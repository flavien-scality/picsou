variable "os" {
  default = "ubuntu-xenial"
}

variable "aws_amis" {
  type = "map"
  default = {
    "ubuntu-xenial" = "ami-ede2e889"
  }
}
