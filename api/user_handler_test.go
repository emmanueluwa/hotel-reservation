package api

import (
    "testing"
    "encoding/json"
    "bytes"

    "net/http/httptest"
    "github.com/gofiber/fiber/v2" 
    "github.com/emmanueluwa/hotel-reservation/types"
)


func TestPostUser(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)

    app := fiber.New()
    userHandler := NewUserHandler(tdb.User)
    app.Post("/", userHandler.HandlePostUser)
    
    params := types.CreateUserParams{
        Email: "grabber@mail.com",
        FirstName: "grab",
        LastName: "muli",
        Password: "pass1234567",
    }

    //marshal params to bytes
    b, _ := json.Marshal(params)

    req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
    req.Header.Add("Content-Type", "application/json")

    resp, err := app.Test(req)
    if err != nil {
        t.Error(err)
    }

    var user types.User
    
    json.NewDecoder(resp.Body).Decode(&user)
    
    if len(user.ID) == 0 {
        t.Errorf("expecting a user id to be set")
    }

    if len(user.EncryptedPassword) > 0 {
        t.Errorf("expecting the Encrypted Password to not be in json response")
    }

    if user.FirstName != params.FirstName {
        t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
    }
    
    if user.LastName != params.LastName {
        t.Errorf("expected last name %s but got %s", params.LastName, user.LastName)
    }

     if user.Email != params.Email {
         t.Errorf("expected email %s but got %s", params.Email, user.Email)
    }
}
