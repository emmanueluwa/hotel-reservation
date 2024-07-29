package api

import (
    "net/http"


	"github.com/gofiber/fiber/v2"
)


func ErrorHandler(c *fiber.Ctx, err error) error {
    if apiError, ok := err.(Error); ok {
        return c.Status(apiError.Code).JSON(apiError)
    }
    //return error from handler but make it internal server error
    apiError := NewError(http.StatusInternalServerError, err.Error())
    return c.Status(apiError.Code).JSON(apiError)
}



type Error struct {
    Code int `json:"code"`
    Err string `json:"error"`
}


//implementing error interface
func (e Error) Error() string {
    return e.Err
}

func NewError(code int, err string) Error {
    return Error{
        Code: code,
        Err: err,
    }
}


func ErrResourceNotFound(res string) Error {
    return Error{
        Code: http.StatusNotFound,
        Err: res + " resource not found",
    }
}


func ErrBadRequest() Error {
    return Error{
        Code: http.StatusBadRequest,
        Err: "invalid JSON request",
    }
}


func ErrUnAuthorized() Error {
    return Error{
        Code: http.StatusUnauthorized,
        Err: "unathorized request",
    }
}


func ErrInvalidID() Error {
    return Error{
        Code: http.StatusBadRequest,
        Err: "invalid id given",
    }
}
