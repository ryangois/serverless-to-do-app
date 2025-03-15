resource "aws_lambda_function" "to-do-lambda" {
  function_name = var.lambda_fucntion_name
  role = aws_iam_role.lambda_role.arn
  runtime = "go1.x"
  handler = "main"
  filename = "lambda.zip"
}