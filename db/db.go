package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBNAME = "hotel-reservation"

func ToObjectID(id string) primitive.ObjectID {
    old, err := primitive.ObjectIDFromHex(id) 
    if err != nil {
        panic(err)
    }
    return old
}
