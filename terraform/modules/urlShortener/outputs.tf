// grab this url to interact with the api_client
output "api_url"{ // url of the api gateway
  value = aws_api_gateway_stage.dev.invoke_url
}
