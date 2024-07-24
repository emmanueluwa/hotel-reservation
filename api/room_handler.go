package api

import (
    "time"
    "fmt"
    "net/http"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/emmanueluwa/hotel-reservation/db"
    "github.com/emmanueluwa/hotel-reservation/types"
    "github.com/gofiber/fiber/v2"
)

type BookRoomResponse struct {
    FromDate time.Time `json:"fromDate"`
    TillDate time.Time `json:"tillDate"`
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

    booking := types.Booking{
        UserID: user.ID,
        RoomID: roomID,
    }

    fmt.Println(booking)
    return nil

}
