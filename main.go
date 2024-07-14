package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/emmanueluwa/hotel-reservation/api"
	"github.com/emmanueluwa/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userCollection = "users"

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	collection := client.Database(dbname).Collection(userCollection)
	
	user := types.User{
		FirstName: "Jack",
		LastName: "101",
	}

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New()

	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)


	//boot up api server
	app.Listen(*listenAddr)
}


