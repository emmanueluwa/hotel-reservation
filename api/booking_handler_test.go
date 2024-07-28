package api

import (
    "testing"
    "fmt"
    "time"
    "net/http/httptest"
    "net/http"
    "encoding/json"

    "github.com/gofiber/fiber/v2" 
    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/emmanueluwa/hotel-reservation/fixtures"
    "github.com/emmanueluwa/hotel-reservation/middleware"
)


func TestUserGetBooking(t *testing.T) {
    db := setup(t)
    defer db.teardown(t)

    var (
        nonAuthUser = fixtures.AddUser(db.Store, "jack", "bread", false)
        user = fixtures.AddUser(db.Store, "jack", "stack", false)
        hotel = fixtures.AddHotel(db.Store, "novotel", "london", 5, nil)
        room = fixtures.AddRoom(db.Store, "couple", true, 9999.9, hotel.ID)
        from = time.Now()
        till = from.AddDate(0, 0, 5)
        booking = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
        app = fiber.New()
        route  = app.Group("/", middleware.JWTAuthentication(db.User))
        bookingHandler = NewBookingHandler(db.Store)
    )
 
    route.Get("/:id", bookingHandler.HandleGetBooking)
    req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("non 200 code got %d", resp.StatusCode)
    }

    var bookingResp *types.Booking

    if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
        t.Fatal(err)
    }

    if bookingResp.ID != booking.ID {
        t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
    }
    
    if bookingResp.UserID != booking.UserID {
        t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
    }

    req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

    resp, err = app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    if resp.StatusCode == http.StatusOK {
        t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
    }


}


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


    // test if not admin user cannot access bookings
    req = httptest.NewRequest("GET", "/", nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
    resp, err = app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    if resp.StatusCode == http.StatusOK {
        t.Fatalf("expected non 200 statusCode but got this %d", resp.StatusCode)
    }

}
