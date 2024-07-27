package main

import (
    "context"
    "log"
    "fmt"
    "time"

    "github.com/emmanueluwa/hotel-reservation/api"
    "github.com/emmanueluwa/hotel-reservation/fixtures"
    "github.com/emmanueluwa/hotel-reservation/db"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

)

func main() {
    ctx := context.Background()
    var err error
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }

    //reset
    if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
        log.Fatal(err)
    }
    

    hotelStore := db.NewMongoHotelStore(client)

    store := &db.Store{
        User: db.NewMongoUserStore(client),
        Booking: db.NewMongoBookingStore(client),
        Room:  db.NewMongoRoomStore(client, hotelStore),
        Hotel: hotelStore,
    }

    user := fixtures.AddUser(store, "jack", "stack", false)
    fmt.Println("jack =>", api.CreateTokenFromUser(user))
    admin := fixtures.AddUser(store, "admin", "admin", true)
    fmt.Println("admin =>", api.CreateTokenFromUser(admin))

    hotel := fixtures.AddHotel(store, "novotel", "london", 5, nil)

    room := fixtures.AddRoom(store, "large", true, 999.99, hotel.ID)

    booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,5))    

    fmt.Println("booking =>", booking.ID)

}


