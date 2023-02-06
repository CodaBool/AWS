resource "aws_iam_role" "cw_assume" {
  name               = "cloudwatch-assume"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_instance_profile" "ec2_profile" {
  name = "ec2_profile"
  role = aws_iam_role.cw_assume.name
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.cw_assume.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "logs" {
  role       = aws_iam_role.cw_assume.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
}

resource "aws_iam_role_policy_attachment" "retention" {
  role       = aws_iam_role.cw_assume.name
  policy_arn = aws_iam_policy.retention.arn
}

resource "aws_iam_policy" "retention" {
  name        = "change_retention"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = [
        "logs:PutRetentionPolicy"
      ]
      Effect   = "Allow"
      Resource = "*"
    }]
  })
}

resource "aws_ssm_parameter" "agent" {
  name  = "agent"
  type  = "SecureString"
  value = file("${path.module}/agent.json")
}
