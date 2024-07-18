package api

import
( 
    "fmt"
    "errors"
   
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/emmanueluwa/hotel-reservation/db"
    
    "github.com/emmanueluwa/hotel-reservation/types"
    
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



func (h *UserHandler) HandlerPutUser(c *fiber.Ctx) error {
    return nil
}


func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
    userID := c.Params("id")
    if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
        return err
    }
    return c.JSON(map[string]string{"deleted": userID})
}


func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
    var params types.CreateUserParams
    if err := c.BodyParser(&params); err != nil {
        return err
    }

    if errors := params.Validate(); len(errors) > 0 {
        return c.JSON(errors)
    }

    user, err := types.NewUserFromParams(params)
    if err != nil {
        return err
    }

    insertedUser, err := h.userStore.InsertUser(c.Context(), user)
    if err != nil {
        return err
    }
    return c.JSON(insertedUser)
}

func (h *UserHandler)  HandleGetUser(c *fiber.Ctx) error {
    var (
        id = c.Params("id")
    )

    user, err := h.userStore.GetUserByID(c.Context(), id)
    if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            return c.JSON(map[string]string{"error": "not found"})
        }
        return err
    }

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
    if err != nil {
        return err
    }

    fmt.Println(users)
    return c.JSON(users)
}


