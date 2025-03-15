resource "aws_iam_role" "lambda_role" {
  name = "lambda_exec_role"

  assume_role_policy = jsondecode({
    Version = "2012-10-17"
    Statement = [{
        Effect = "Allow"
        Principal = { Service = "lambda.amazonaws,com"}
        Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_policy_attachment" "lambda_dynamodb" {
  name = "lambda_dynamodb_policy"
  roles = [ aws_iam_role.lambda_role.name ]
  policy_arn = "arn:aws:iam::aws:policy/AmazonDyunamoDBFullAccess"
}