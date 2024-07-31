package main

import (
	"context"
	"log"
    "os"

	"github.com/emmanueluwa/hotel-reservation/api"
	"github.com/emmanueluwa/hotel-reservation/db"
    "github.com/joho/godotenv" 
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
   )


// Configuration
// 1. MONGODB endpoitn
// 2. ListenAddress for HTTP server
// 3. JWT secret
// 4. MONGODBNAME

var config = fiber.Config{
    ErrorHandler: api.ErrorHandler,
}

func main() {
    mongoEndpoint := os.Getenv("MONGO_DB_URL")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI( mongoEndpoint ))
    if err != nil {
        log.Fatal(err)
    }

    var (
        //initialising handlers
        hotelStore = db.NewMongoHotelStore(client)
        roomStore = db.NewMongoRoomStore(client, hotelStore)
        userStore = db.NewMongoUserStore(client)
        bookingStore = db.NewMongoBookingStore(client)

        store = &db.Store{
            Hotel: hotelStore,
            Room: roomStore,
            User: userStore,
            Booking: bookingStore,
        }

        userHandler = api.NewUserHandler(userStore)
        hotelHandler = api.NewHotelHandler(store)
        authHandler = api.NewAuthHandler(userStore)
        roomHandler = api.NewRoomHandler(store)
        bookingHandler = api.NewBookingHandler(store)
        
        app = fiber.New(config)
        auth = app.Group("/api")

	    apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
        
        admin = apiv1.Group("/admin", api.AdminAuth)
    )


    //auth
    auth.Post("/auth", authHandler.HandleAuthenticate)
	
    //VERSIONED API ROUTES
    //user handlers
    apiv1.Post("/user", userHandler.HandlePostUser)    
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
    
    //hotel handlers
    apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
    apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
    apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

    // room handlers
    apiv1.Get("/room", roomHandler.HandleGetRooms)
    apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
    //TODO: cancel booking

    //booking handlers
    admin.Get("/booking", bookingHandler.HandleGetBookings)

    apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
    apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

    //boot up api server
    listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}


func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal(err)
    }  
}
