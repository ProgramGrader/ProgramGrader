terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 4.10"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}

resource "aws_s3_bucket" "url_s3_b" {
  bucket = "url_s3_bucket"

  tags = {
    Name = "url_s3_bucket"
    Environment = "Dev"
  }

}

resource "aws_dynamodb_table" "urls" {
  name = "urls"
  billing_mode = "PROVISIONED"
  attribute {
    name = "urlId"
    type = "S"
  }
  hash_key = "urlId"
}

module "table_autoscaling"  {
  source = "snowplow-devops/dynamodb-autoscaling/aws"
  table_name = aws_dynamodb_table.urls
}


// url lambda


