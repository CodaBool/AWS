output "role" {
  description = "Role name to attach policy to"
  value       = aws_iam_role.lambda_assume.name
}

output "function" {
  value       = aws_lambda_function.main
}