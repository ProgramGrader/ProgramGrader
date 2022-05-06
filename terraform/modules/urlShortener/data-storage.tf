
resource "aws_s3_bucket" "url_s3_b" {
  bucket = "url-s3-bucket"
}

resource "aws_dynamodb_table" "urls" {
  name = "S3URLS"
  billing_mode = "PAY_PER_REQUEST"
  attribute {
    name = "urlId"
    type = "S"
  }
  hash_key = "urlId"
}
