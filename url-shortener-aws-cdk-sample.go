package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type UrlShortenerAwsCdkSampleStackProps struct {
	awscdk.StackProps
}

func NewUrlShortenerAwsCdkSampleStack(scope constructs.Construct, id string, props *UrlShortenerAwsCdkSampleStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	table := awsdynamodb.NewTable(stack, jsii.String("mapping-table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("Id"), Type: awsdynamodb.AttributeType_STRING},
	})

	lambdaAsset := awss3assets.NewAsset(stack, jsii.String("lambda-function-zip"), &awss3assets.AssetProps{
		Path: jsii.String("./functions/main.zip"),
	})

	function := awslambda.NewFunction(stack, jsii.String("back-end"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("main"),
		Code:    awslambda.Code_FromBucket(lambdaAsset.Bucket(), lambdaAsset.S3ObjectKey(), nil),
	})

	table.GrantReadWriteData(function)
	function.AddEnvironment(jsii.String("TABLE_NAME"), table.TableName(), nil)

	awsapigateway.NewLambdaRestApi(stack, jsii.String("api"), &awsapigateway.LambdaRestApiProps{
		Handler: function,
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewUrlShortenerAwsCdkSampleStack(app, "UrlShortenerAwsCdkSampleStack", &UrlShortenerAwsCdkSampleStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
