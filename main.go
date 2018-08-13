package main

import (
	"fmt"
	"github.com/falafeljan/gin-simple-token-middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"strings"
)

func main() {
	args := getArgs()
	options := redis.Options{
		Addr: strings.Join([]string{args.RedisHost, args.RedisPort}, ":"),
		DB:   0,
	}

	client := redis.NewClient(&options)
	tokenContext := NewTokenContext(client, args.RedisPrefix)

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	if args.InProduction {
		gin.SetMode("release")
	}

	token, err := generateRandomString(32)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The security token for this instance is:\n\t%s\n"+
		"Please use this token with the 'Authorization' HTTP header as `Token <token>`.\n"+
		"For more information, please refer to the documentation.\n\n", token)

	router := gin.Default()
	router.Use(createCORSMiddleware(args.AllowedOrigins))
	router.Use(limiterMiddleware(client, args.RateLimit))

	group := router.Group("/", tokenmiddleware.NewHandler(token))
	group.POST("/graphql", createGraphQLHandler(tokenContext, args))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.HTTPPort}, "")))
}
