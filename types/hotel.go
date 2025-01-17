package types

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name string `bson:"name" json:"name"`
    Location string `bson:"location" json:"location"`
    Rooms []primitive.ObjectID `bson:"rooms" json:"rooms"`
    Rating int `bson:"rating" json:"rating"`
}

type Room struct {
    ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    // single, family, couple
    Size string `bson:"size" json:"size"`
    Penthouse bool `bson:"penthouse" json:"penthouse"`
    Price float64 `bason:"price" json:"price"`
    HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

