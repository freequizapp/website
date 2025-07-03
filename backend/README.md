# use the AWS SAM CLI for testing
### Build the container using SAM
`sam build`

### We have to test after deployment, for streaming:
`sam sync --watch \
  --stack-name freequizapp \
  --parameter-overrides GROQAPIKEY=<insert-api-key>\
  --region us-west-2
`

### test with command line request
`curl -iN -X POST https://<insert-lambda-url> \
  -H "Content-Type: application/json" \
  -d '{"prompt": "javascript"}'
`

### For deployment we can use something like this
`sam deploy --parameter-overrides GROQAPIKEY={github.secrets.groqapikey}`

