package api

import (
    "github.com/emmanueluwa/hotel-reservation/db"
    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
    store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler{
    return &BookingHandler{
        store: store,
    }
}

// this needs to be admin authorised
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
    bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
    if err != nil {
        return err
    }
    return c.JSON(bookings)
}


// this needs to be user authorised
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
    id := c.Params("id")
    booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
    if err != nil {
        return err
    }
    return c.JSON(booking)
}
