// Recipes API

// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version 0.0.1
// Consumes:
// - application/json
// Produces
// - application/json
// swagger:meta
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"recipe-api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client

var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func main() {
	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.GET("/recipes/:id", recipesHandler.GetRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		authorized.GET("/recipes/search", recipesHandler.SearchRecipesHandler)
	}
	router.Run(":8080")
}

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database("demo").Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = &handlers.AuthHandler{}
}
