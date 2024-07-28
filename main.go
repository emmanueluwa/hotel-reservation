package main

import (
	"context"
	"flag"
    "time"	
	"log"
    "net/http"
    "fmt"

    

    "github.com/emmanueluwa/hotel-reservation/middleware"    
	"github.com/emmanueluwa/hotel-reservation/api"
	"github.com/emmanueluwa/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
   )


var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        if apiError, ok := err.(api.Error); ok {
            return c.Status(apiError.Code).JSON(apiError)
        }
        //return error from handler but make it internal server error
        return api.NewError(http.StatusInternalServerError, err.Error())
    },
}

func main() {
    now := time.Now()
    fmt.Println(now)  

    listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
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

	    apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
        
        admin = apiv1.Group("/admin", middleware.AdminAuth)
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
	app.Listen(*listenAddr)
}


