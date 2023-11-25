terraform {
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "terraform-backend-hk"
    key    = "terraform.tfstate"
  }
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}
