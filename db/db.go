package db

import (
)

const MongoDBNameEnv = "MONGO_DB_NAME"

type Pagination struct {
    Limit int64
    Page int64
}

type Store struct {
    User UserStore
    Hotel HotelStore
    Room RoomStore
    Booking BookingStore
}

