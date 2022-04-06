
terraform {
  backend "s3" {
    bucket = "tfstate-3ea6z45i"
    key    = "programGrader/terraform.tfstate"
    dynamodb_table = "terraform_state_lock"
    region = "us-east-2"
  }
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.8"
    }
  }

  required_version = ">= 1.1.0"
}

locals {
  ProjectName = "AutoGrader"
}

provider "aws" {
  profile = "default"
  region  = "us-east-2"

  default_tags {
    tags = {
      Terraform   = "true"
      Project     = local.ProjectName
    }
  }
}

resource "aws_dynamodb_table" "analytics" {
  name         = "AutoGrader_Analytics"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "CourseNumber-SemesterID"
  range_key = "AssignmentName-UserName"

  point_in_time_recovery {
    enabled = true
  }

  attribute {
    name = "CourseNumber-SemesterID"
    type = "S"
  }

  attribute {
    name = "AssignmentName-UserName"
    type = "S"
  }

  lifecycle {
    prevent_destroy = true
  }
}

module "backup" {
  source  = "lgallard/backup/aws"
  version = "0.12.1"
  # insert the 12 required variables here

  plan_name = "Weekly-Monthly-Forever-Retention"

  #vault_name = "Default"

  # Multiple rules using a list of maps
  rules = [
    {
      name                     = "Weekly"
      schedule                 = "cron(0 7 ? * fri *)"
      target_vault_name        = "Default"
      start_window             = 120
      completion_window        = 360
      lifecycle = {
        cold_storage_after = 7
        delete_after = 18262
      }
      copy_action         = {}
      recovery_point_tags = {}
    },
    {
      name                = "Monthly"
      target_vault_name   = "Default"
      schedule            = "cron(0 7 1 * ? *)"
      start_window        = 120
      completion_window   = 360
      lifecycle = {
        cold_storage_after = 30
        delete_after = 18262
      }
      copy_action         = {}
      recovery_point_tags = {}
    },
  ]

  selections = [
    {
      name = local.ProjectName
      selection_tags = [
        {
          type  = "STRINGEQUALS"
          key   = "Project"
          value = local.ProjectName
        }
      ]
    }

  ]

}

resource "aws_s3_bucket" "b1" {
  bucket = format("%s-bucket", lower(local.ProjectName))
  acl    = "private"


  lifecycle {
    prevent_destroy = true
  }

}

resource "aws_s3_bucket_intelligent_tiering_configuration" "intelligent_tier_resource" {
  bucket = aws_s3_bucket.b1.bucket
  name = format("entire_bucket_%s",local.ProjectName)

  tiering {
    access_tier = "DEEP_ARCHIVE_ACCESS"
    days        = 180
  }
  tiering {
    access_tier = "ARCHIVE_ACCESS"
    days        = 90
  }

}
