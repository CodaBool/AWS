resource "aws_s3_bucket" "bucket" {
  bucket = "codabool-go"
  tags = {
    "Name" = "my bucket"
  }
}