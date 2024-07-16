package api

import (
    "context"

	"github.com/emmanueluwa/hotel-reservation/types"
    
	"github.com/emmanueluwa/hotel-reservation/db"
    
	"github.com/gofiber/fiber/v2"

    
)

type UserHandler struct {
        userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
    return &UserHandler{
        userStore: userStore,
    }
}

func (h *UserHandler)  HandleGetUser(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
        ctx = context.Background()
    )

    user, err := h.userStore.GetUserByID(ctx, id)
    if err != nil {
        return err
    }

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Jack",
		LastName: "101",
	}
	return c.JSON(u)
}


