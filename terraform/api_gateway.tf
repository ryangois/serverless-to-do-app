resource "aws_apigatewayv2_api" "http_api" {
  name = var.api_gateway_name
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id = aws_apigatewayv2_api.http_api.id
  name = "$dafault"
  auto_deploy = true
}