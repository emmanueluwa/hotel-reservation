package api

import (
    "testing"
    "time"
    "net/http/httptest"
    "net/http"
    "encoding/json"

    "github.com/gofiber/fiber/v2" 
    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/emmanueluwa/hotel-reservation/fixtures"
    "github.com/emmanueluwa/hotel-reservation/middleware"
)


func TestAdminGetBookings(t *testing.T) {
    db := setup(t)
    defer db.teardown(t)

    var (
        adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
        user = fixtures.AddUser(db.Store, "jack", "stack", false)
        hotel = fixtures.AddHotel(db.Store, "novotel", "london", 5, nil)
        room = fixtures.AddRoom(db.Store, "couple", true, 9999.9, hotel.ID)
        from = time.Now()
        till = from.AddDate(0,0,5)
        booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
        app = fiber.New()
        admin = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
        bookingHandler = NewBookingHandler(db.Store)
    )
    
    admin.Get("/", bookingHandler.HandleGetBookings)
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("non 200 response %d", resp.StatusCode)
    }

    var bookings []*types.Booking
    if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
        t.Fatal(err)
    }
    
    if len(bookings) != 1 {
        t.Fatalf("expected 1 booking got %d", len(bookings))
    }
    
    have := bookings[0]
    if have.ID != booking.ID {
        t.Fatalf("expected %s got %s", booking.ID, have.ID)
    }
    
    if have.UserID != booking.UserID {
        t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
    }
}
