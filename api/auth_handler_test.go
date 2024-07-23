package api

import (
    "testing"
    "fmt"
    "encoding/json"
    "bytes"
    "net/http"
    "net/http/httptest"
    "context"
    "reflect"

    "github.com/gofiber/fiber/v2" 
    "github.com/emmanueluwa/hotel-reservation/db"
    "github.com/emmanueluwa/hotel-reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
    user, err := types.NewUserFromParams(types.CreateUserParams{
        Email: "jack@bread.com",
        FirstName: "jack",
        LastName: "stack",
        Password: "mostsecurepassword",
    })
    if err != nil {
        t.Fatal(err)
    } 

    _, err = userStore.InsertUser(context.TODO(), user)
    if err != nil {
        t.Fatal(err)
    }
    return user
}


func TestAuthenticateWithWrongPassword(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)
    insertTestUser(t, tdb.UserStore)

    app := fiber.New()
    authHandler := NewAuthHandler(tdb.UserStore)
    app.Post("/auth", authHandler.HandleAuthenticate)
    
    params := AuthParams{
        Email: "jack@bread.com",
        Password: "incorrectpassword",
    }

    b, _ := json.Marshal(params)
    req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))

    req.Header.Add("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    if resp.StatusCode != http.StatusBadRequest {
        t.Fatalf("expected http status of 400 but go %d", resp.StatusCode)
    }

    var genResp genericResp
    if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
    t.Fatal(err)
    }

    if genResp.Type != "error" {
        t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
    }

    if genResp.Msg != "invalid credentials" {
        t.Fatalf("expected gen response msg to be <invalid credentials> but got %s", genResp.Msg)
    }
    }


func TestAuthenticateSuccess(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)
    insertedUser := insertTestUser(t, tdb.UserStore)

    app := fiber.New()
    authHandler := NewAuthHandler(tdb.UserStore)
    app.Post("/auth", authHandler.HandleAuthenticate)
    
    params := AuthParams{
        Email: "jack@bread.com",
        Password: "mostsecurepassword",
    }

    b, _ := json.Marshal(params)
    req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))

    req.Header.Add("Content-Type", "application/json")
    
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected http status of 200 but go %d", resp.StatusCode)
    }

    var authResp AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        t.Fatal(err)
    }
    
    if authResp.Token == "" {
        t.Fatalf("expected the JWT token to be present in the auth response")
    
    fmt.Println(insertedUser, authResp.User)

    //setting encrypted password to empty string since it is not returned in any json response
    insertedUser.EncryptedPassword = ""

    if !reflect.DeepEqual(insertedUser, authResp.User) {
        t.Fatalf("expected the user to be the inserted user")
        fmt.Println(insertedUser, authResp.User)
    }

    }
}
