# use the AWS SAM CLI for testing
### Build the container using SAM
`sam build`

### Run it as a local http sever
`sam local start-api --port 8080 --env-vars env.json`

### Hit the api
`
curl -X POST http://localhost:8080/generate-questions \
  -H "Content-Type: application/json" \
  -d '{"prompt": "rocket engines"}'
`

Or use the frontend to prompt.
