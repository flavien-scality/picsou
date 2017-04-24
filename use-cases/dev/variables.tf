variable "os" {
  type = "string"
  default = "ubuntu-xenial"
}

variable "login" {
  type = "string"
  default = "hsimpson"
}

# AWS specific

variable "vpc_id" {
  type = "string"
  description = "VPC ID for AWS resources."
}

variable "availability_zone_id" {
  type = "string"
  description = "AZ used to create EC2 instances."
}

variable "subnet_id" {
  type = "string"
  description = "Subnet for EC2 instances."
}
