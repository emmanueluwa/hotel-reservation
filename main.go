package main

import (
	"context"
	"flag"
	
	"log"


	"github.com/emmanueluwa/hotel-reservation/api"
	"github.com/emmanueluwa/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
   )


const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userCollection = "users"

var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
    	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
    if err != nil {
        log.Fatal(err)
    }

    //initialising handlers
    userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

    apiv1.Post("/user", userHandler.HandlePostUser)    
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)


	//boot up api server
	app.Listen(*listenAddr)
}


