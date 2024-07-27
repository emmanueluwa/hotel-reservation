package api

import (
    "testing"
    "time"
    "fmt"

    "github.com/emmanueluwa/hotel-reservation/fixtures"
)


func TestAdminGetBookings(t *testing.T) {
    db := setup(t)
    defer db.teardown(t)

    user := fixtures.AddUser(db.Store, "jack", "stack", false)
    hotel := fixtures.AddHotel(db.Store, "novotel", "london", 5, nil)
    room := fixtures.AddRoom(db.Store, "couple", true, 9999.9, hotel.ID)

    from := time.Now()
    till := from.AddDate(0,0,5)
    booking := fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)

    fmt.Println(booking)
}
