package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

var (
	client   *redis.Client
	headers  map[string]string
	redisURL string
)

const (
	redisURLEnv      = "REDIS"
	notFoundResponse = `<!DOCTYPE html>
<html>
	<body>
		<h3>link not found, try another code</h3>
	</body>
</html>
`
	redirectResponse = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="refresh" content="0; url=%s">
	</head>
</html>
`
)

func init() {
	redisURL = os.Getenv(redisURLEnv)

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Println(err)
	}
	client = redis.NewClient(opt)

	headers = make(map[string]string)
	headers["Content-Type"] = "text/html"
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

	hash := request.Headers["hash"]
	val, err := client.Get(ctx, hash).Result()
	if err == redis.Nil {
		return Response{
			StatusCode: 404,
			Body:       notFoundResponse,
			Headers:    headers,
		}, nil
	} else if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: 200,
		Body:       fmt.Sprintf(redirectResponse, val),
		Headers:    headers,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
