package main

import (
    "context"
    "log"

    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/emmanueluwa/hotel-reservation/db"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"

)

var (
    client *mongo.Client
    roomStore db.RoomStore
    hotelStore db.HotelStore
    ctx = context.Background()
)

func seedHotel(name, location string) {
     hotel := types.Hotel{
        Name: name,
        Location: location,
        Rooms: []primitive.ObjectID{},
    }

    rooms := []types.Room{
        {
            Type: types.PenthouseRoomType,
            BasePrice: 200.0,    
        },
        {
            Type: types.DeluxeRoomType,
            BasePrice: 450.0,
        },
        {
            Type: types.CouplesRoomType,
            BasePrice: 300.0,
        },
    }

    insertedHotel, err := hotelStore.Insert(ctx, &hotel)
    if err != nil {
        log.Fatal(err)
    }

    for _, room := range rooms{
        room.HotelID = insertedHotel.ID
        _, err := roomStore.InsertRoom(ctx, &room)
        if err != nil {
            log.Fatal(err)
        }
    }
}

func main() {
    seedHotel("Novotel", "London")
    seedHotel("Hotel Colline de France", "Brazil") 
    seedHotel("Londolozi Game reserve", "South Africa")
}


func init() {
    var err error
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }

    //reset
    if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
        log.Fatal(err)
    }

    hotelStore = db.NewMongoHotelStore(client)
    roomStore = db.NewMongoRoomStore(client, hotelStore)


}
