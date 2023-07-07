output "instance" {
  value = aws_spot_instance_request.main
}

output "eip" {
  value = aws_eip.main
}