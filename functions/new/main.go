package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Input struct {
	Link string `json:"link"`
}

var (
	client   *redis.Client
	redisURL string
)

const (
	redisURLEnv = "REDIS"
)

func init() {

	redisURL = os.Getenv(redisURLEnv)
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Println(err)
	}
	client = redis.NewClient(opt)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	b := request.Body

	var i Input
	err := json.Unmarshal([]byte(b), &i)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	h := hash(i.Link)[:5]
	log.Println(h)

	_, err = client.Get(ctx, h).Result()

	// no value for key
	if err == redis.Nil {
		err = client.Set(ctx, h, i.Link, 0).Err()
		if err != nil {
			return Response{}, err
		}

		return Response{
			StatusCode: 200,
			Body:       h,
		}, nil

	} else if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: 200,
		Body:       h,
	}, nil

}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	lambda.Start(Handler)
}
