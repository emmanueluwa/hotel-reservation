package api

import (
    "github.com/gofiber/fiber/v2"
    "github.com/emmanueluwa/hotel-reservation/types"
)

func AdminAuth(c *fiber.Ctx) error {
    user, ok := c.Context().UserValue("user").(*types.User)
    if !ok {
        return ErrUnAuthorized()   
    }

    if !user.IsAdmin{
        return ErrUnAuthorized()
    }
    return c.Next()
}
