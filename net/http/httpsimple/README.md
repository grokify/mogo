# HTTP Simple for Go

## HTTP Simple Server

### AWS Lambda Manual Setup

* AWS Lambda Function Configuration
  1. Create AWS Lambda function
  2. Set `SimpleServer.HttpEngine` to `awslambda`, e.g. using environment variable
  4. Set Handler to `main`
* API Gateway Configuration
  1. Create "REST" API
  2. Select "New API"
  3. Select "Actions" > "Create Resource"
  4. Click "Configure as proxy resource"
  5. Use Resource Path `{proxy+}`
  6. Click "Enable API Gateway CORS"
  7. Click "Create Resource"
  8. Leave "Integration Type" as "Lambda Function Proxy"
  9. In "Lambda Function", paste in your Lamda ARN
 10. Click "Deploy API" and create stage if necessary
 11. Test with cURL against Stage Invoke URL
* API Gateway API Key
  1. Create API Key
  2. Create Usage Plan
  3. Add API Stage by Selecting API Gateway, Stage, and clicking checkmark icon.
  4. Click "Next".
  5. Add API Key, click checkmark icon, click "Done"
  6. Click on API Gateway `{proxy+}` endpoint, click on `ANY`
  7. Click on "Method Request"
  8. Click edit icon by "API Key Required", change to true, click checkmark icon
  9. Select "Deploy API"
 10. Test with cURL and `X-API-Key` header

