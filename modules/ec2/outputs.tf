# output "instance" {
#   value = aws_spot_instance_request.main
# }
output "instance" {
  value = aws_instance.main
}