package main

import (
    "fmt"
    "context"
    "log"

    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/emmanueluwa/hotel-reservation/db"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
     "go.mongodb.org/mongo-driver/bson/primitive"   
)

func main() {
    ctx := context.Background()

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }

    //reset
    if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
        log.Fatal(err)
    }

    hotelStore := db.NewMongoHotelStore(client)
    roomStore := db.NewMongoRoomStore(client, hotelStore)

    hotel := types.Hotel{
        Name: "Novotel",
        Location: "London",
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


    insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
    if err != nil {
        log.Fatal(err)
    }

    for _, room := range rooms{
        room.HotelID = insertedHotel.ID
        insertedRoom, err := roomStore.InsertRoom(ctx, &room)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(insertedRoom)
    }

}
