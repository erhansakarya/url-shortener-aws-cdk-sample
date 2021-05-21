module url-shortener-aws-cdk-sample

go 1.16

require (
	github.com/aws/aws-cdk-go/awscdk v1.104.0-devpreview
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go-v2 v1.6.0
	github.com/aws/aws-sdk-go-v2/config v1.3.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.1.1
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.3.1
	github.com/aws/constructs-go/constructs/v3 v3.3.75
	github.com/aws/jsii-runtime-go v1.29.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.10.0
	github.com/gin-gonic/gin v1.7.1
	github.com/google/uuid v1.2.0
	github.com/stretchr/testify v1.7.0

	// for testing
	github.com/tidwall/gjson v1.7.4
)
