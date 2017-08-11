provider "aws" {
  region     = "eu-west-2"
  shared_credentials_file = "${pathexpand("~/.aws/credentials")}"
}
