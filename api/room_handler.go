package api

import (
    "time"
    "net/http"
    "fmt"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/emmanueluwa/hotel-reservation/db"
    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/gofiber/fiber/v2"
)

type BookRoomParams struct {
    FromDate time.Time `json:"fromDate"`
    TillDate time.Time `json:"tillDate"`
    NumPersons int `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
    now := time.Now()
    if now.After(p.FromDate) || now.After(p.TillDate) {
        return fmt.Errorf("please choose a future date!")
    }
    return nil
}

type RoomHandler struct {
    store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
    return &RoomHandler{
        store: store,
    }
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
    var params BookRoomParams
    if err := c.BodyParser(&params); err != nil {
        return err
    }

    if err := params.validate(); err != nil {
        return err
    }

    roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return err
    }
    
    user, ok := c.Context().Value("user").(*types.User)
    if !ok {
        return c.Status(http.StatusInternalServerError).JSON(genericResp{
            Type: "error",
            Msg: "internal server error",
        })
    }

    where := bson.M{
        "fromDate": bson.M{
            "$gte": params.FromDate,
        },
        "tillDate": bson.M{
            "$lte": params.TillDate,
        },
    }

    bookings, err := h.store.Booking.GetBookings(c.Context(), where)
    if err != nil {
        return err
    }

    if len(bookings) > 0 {
        return c.Status(http.StatusBadRequest).JSON(genericResp{
            Type: "error",
            Msg: fmt.Sprintf("room %s already booked", c.Params("id")),
        })
    }

    booking := types.Booking{
        UserID: user.ID,
        RoomID: roomID,
        FromDate: params.FromDate,
        TillDate: params.TillDate,
        NumPersons: params.NumPersons,
    }

    inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
    if err != nil {
        return err
    }

    return c.JSON(inserted)
}
