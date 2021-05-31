package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

func getClient() redis.Client {
	opt, _ := redis.ParseURL("rediss://:c9b58dc0c64e41bd9de152d2359de387@eu1-holy-kiwi-32130.upstash.io:32130")
	return *redis.NewClient(opt)
}

func main() {
	client := getClient()

	ctx := context.Background()
	err := client.Set(ctx, "hello", "world", 0).Err()
	if err != nil {
		log.Println(err)
	}

	val, err := client.Get(ctx, "golang").Result()
	if err == redis.Nil {
		log.Println("no such key")
	} else if err != nil {
		log.Println(err)
	}

	fmt.Println(val)

	err = client.Del(ctx, "hello").Err()
	if err != nil {
		log.Println()
	}

	var (
		results []string
		scan []string
		cursor uint64
	)

	for {
		scan, cursor, err = client.Scan(context.Background(), cursor, "", 20).Result()
		if err != nil {
			log.Println(err)
		}

		results = append(results, scan...)

		if cursor == 0 {
			break
		}
	}

	for _, item := range results {
		err = client.Del(ctx, item).Err()
		if err != nil {
			log.Println()
		}
	}
}