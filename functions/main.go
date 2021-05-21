package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ginLambda *ginadapter.GinLambda
var cfg aws.Config

// Request model definition
type UrlCreationRequest struct {
	LongUrl string `json:"long_url"`
}

// init the Gin Server
func init() {
	var configErr error
	cfg, configErr = config.LoadDefaultConfig(context.TODO())
	if configErr != nil {
		log.Fatalf("unable to load SDK config, %v", configErr)
	}

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the URL Shortener API",
		})
	})
	r.POST("/create-short-url", createShortURL)
	r.GET("/:shortUrl", redirectShortURL)

	ginLambda = ginadapter.New(r)
}

type Item struct {
	Id        string
	TargetURL string
}

func createShortURL(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableName := os.Getenv("TABLE_NAME")

	targetURL := creationRequest.LongUrl

	id := uuid.New().String()[:8]

	dynamoClient := dynamodb.NewFromConfig(cfg)

	item := Item{
		Id:        id,
		TargetURL: targetURL,
	}

	dynamoClient.PutItem(context.Background(),
		&dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]types.AttributeValue{
				"Id":        &types.AttributeValueMemberS{Value: item.Id},
				"TargetURL": &types.AttributeValueMemberS{Value: item.TargetURL},
			},
		})

	url := "https://" + c.Request.URL.Hostname() + "/prod/" + id

	c.JSON(http.StatusOK, gin.H{
		"message": url,
	})
}

func redirectShortURL(c *gin.Context) {
	shortUrlParam := c.Param("shortUrl")
	log.Printf("shortUrlParam is %v", shortUrlParam)

	tableName := os.Getenv("TABLE_NAME")

	dynamoClient := dynamodb.NewFromConfig(cfg)

	response, err := dynamoClient.GetItem(context.Background(),
		&dynamodb.GetItemInput{
			TableName: &tableName,
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberS{Value: shortUrlParam},
			},
			ProjectionExpression: aws.String("TargetURL"),
		})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if response.Item == nil {
		log.Fatalf("item not found")
	}

	item := struct {
		TargetURL string
	}{}

	err = attributevalue.UnmarshalMap(response.Item, &item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("item.TargetURL is %v", item.TargetURL)
	log.Printf("response.item[TargetURL] is %v", response.Item["TargetURL"])
	c.Redirect(302, item.TargetURL)
}

// Handler will deal with Gin working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
