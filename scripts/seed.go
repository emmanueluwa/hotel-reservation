package main

import (
    "context"
    "log"
    "fmt"
    "time"

    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/emmanueluwa/hotel-reservation/api"
    "github.com/emmanueluwa/hotel-reservation/fixtures"
    "github.com/emmanueluwa/hotel-reservation/db"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"

)

var (
    client *mongo.Client
    roomStore db.RoomStore
    hotelStore db.HotelStore
    userStore db.UserStore
    bookingStore db.BookingStore
    ctx = context.Background()
)


func seedUser(isAdmin bool, fname, lname, email string, password string) *types.User {
    user, err := types.NewUserFromParams(types.CreateUserParams{
        Email: email,
        FirstName: fname,
        LastName: lname,
        Password: password, 

    })
    if err != nil {
        log.Fatal(err)
    } 

    user.IsAdmin = isAdmin

    insertedUser, err := userStore.InsertUser(context.TODO(), user)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
    return insertedUser
}


func seedRoom(size string, pent bool, price float64, hotelID primitive.ObjectID) *types.Room {

    room := &types.Room{
        Size: size,
        Penthouse: pent,
        Price: price,
        HotelID: hotelID,
    } 
    
    insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
    if err != nil {
        log.Fatal(err)
    }

    return insertedRoom
}


func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
    booking := &types.Booking{
        UserID: userID,
        RoomID: roomID,
        FromDate: from,
        TillDate: till,
    }

    resp, err := bookingStore.InsertBooking(context.Background(), booking); 
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("resp =>", resp.ID)
}


func seedHotel(name string, location string, rating int) *types.Hotel {
    hotel := types.Hotel{
        Name: name,
        Location: location,
        Rooms: []primitive.ObjectID{},
        Rating: rating,
    }

    insertedHotel, err := hotelStore.Insert(ctx, &hotel)
    if err != nil {
        log.Fatal(err)
    }
    return insertedHotel
}

func main() {
        var err error
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
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
    
    hotel := fixtures.AddHotel(store, "novotel", "london", 5, nil)

    room := fixtures.AddRoom(store, "large", true, 999.99, hotel.ID)

    booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0,0,5))    
    fmt.Println(booking)
    return 

  //  jack := seedUser(false, "jack", "stack", "jack@bread.com", "mostsecurepassword")
//    seedUser(true, "admin", "admin", "admin@bread.com", "mostadminpassword")

  //  seedHotel("Novotel", "London", 5)
    //seedHotel("Hotel Colline de France", "Brazil", 5) 
    //hotel = seedHotel("Londolozi Game reserve", "South Africa", 5)
    //seedRoom("couple", true, 10000, hotel.ID)
    //seedRoom("family", true, 15000, hotel.ID)
    //room := seedRoom("single", true, 8000, hotel.ID)
    //seedBooking(jack.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
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
    userStore = db.NewMongoUserStore(client)
    bookingStore = db.NewMongoBookingStore(client)


}
