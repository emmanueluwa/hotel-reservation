package api

import (
    
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/emmanueluwa/hotel-reservation/db"
    "github.com/gofiber/fiber/v2"
)


type HotelHandler struct {
    hotelStore db.HotelStore
    roomStore db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
    return &HotelHandler{
        hotelStore: hs,
        roomStore: rs,
    }
}


func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
    id := c.Params("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    filter := bson.M{"hotelID": oid}
    rooms, err := h.roomStore.GetRooms(c.Context(), filter)
    if err != nil {
        return err
    }

    return c.JSON(rooms)
}


func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
    hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
    if err != nil {
        return err
    }
    return c.JSON(hotels)
}
