package api

import (
    "net/http"
)

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


func ErrUnAuthorised() Error {
    return Error{
        Code: http.StatusUnauthorized,
        Err: "unathorised request",
    }
}


func ErrInvalidID() Error {
    return Error{
        Code: http.StatusBadRequest,
        Err: "invalid id given",
    }
}
