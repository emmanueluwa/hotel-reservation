package main

import (
    "context"
    "log"
    "math/rand"
    "fmt"
    "os"
    "time"

    "github.com/joho/godotenv"
    "github.com/emmanueluwa/hotel-reservation/api"
    "github.com/emmanueluwa/hotel-reservation/fixtures"
    "github.com/emmanueluwa/hotel-reservation/db"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal(err)
    }

    var (
        ctx = context.Background()
        mongoEndpoint = os.Getenv("MONGO_DB_URL")      
        mongoDBName = os.Getenv("MONGO_DB_NAME")    
    )

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
    if err != nil {
        log.Fatal(err)
    }

    //reset
    if err := client.Database(mongoDBName).Drop(ctx); err != nil {
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

    booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))    

    fmt.Println("booking =>", booking.ID)

    for i := 0; i < 100; i++ {
        name := fmt.Sprintf("sakinah hotel %d", i)
        location := fmt.Sprintf("location %d", i)
        fixtures.AddHotel(store, name, location,  rand.Intn(5)+1, nil)
    }

}


