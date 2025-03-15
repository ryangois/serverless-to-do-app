variable "aws_region" {
  default = "us-east-1"
}

variable "dynamo_table_name" {
  default = "to-dos"
}

variable "lambda_fucntion_name" {
    default = "to-do-lambda"
}

variable "api_gateway_name" {
  default = "to-do-api"
}