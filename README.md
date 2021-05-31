# serverless-link-shortener
![build](https://github.com/mathisve/serverless-link-shortener//actions/workflows/go.yaml/badge.svg?branch=master)

**A stateful serverless API built on top of AWS Lambda and Upstash Redis.**

## Special thanks to [Upstash](https://upstash.com/) for making projects like this possible!

## About:

The codebase consists of 2 functions:
- `new` to add links to the database, returns the shortened hash/code
- `get` to retrieve the links given the hash/code, returns HTML

## Code snippet:
```go
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
    // extract hash from API Gateway headers
    hash := request.Headers["hash"]
    
    // get from Redis compatible Upstash database
    val, err := client.Get(ctx, hash).Result()
    
    if err == redis.Nil {
        // hash not found
        return Response{
            StatusCode: 404,
            Body:       notFoundResponse,
            Headers:    headers,
        }, nil
    }

    // did find hash
    return Response{
        StatusCode: 200,
        Body:       fmt.Sprintf(redirectResponse, val),
        Headers:    headers,
    }, nil
}
```

## Video:
If you prefer to have me explain it to you, you're in luck!
Watch the video [here](https://youtu.be/EJ6CJ0GC9lk)! (or click the thumbnail!)
[![](https://raw.githubusercontent.com/mathisve/serverless-link-shortener/master/img/upstashthumb.png?token=ACUQMQPCQQROUJ3NS64FZE3AWVU6M)](https://youtu.be/EJ6CJ0GC9lk)
