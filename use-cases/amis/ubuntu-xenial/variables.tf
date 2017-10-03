variable "region" {
  type = "string"
  default = "eu-west-2"
}

variable "instance_type" {
  type = "string"
  default = "t2.micro"
}

variable "key_name" {
  type = "string"
  default = "maxime-london"
}

variable "filters" {
  type = "map"
  default = {
    name   = "name"
    values = ["testing-*"]
  }
}
