package types

import (
    "fmt"
    "regexp"
    "golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const (
    bcryptCost = 12
    minFirstNameLen = 2
    minLastNameLen = 2
    minPasswordLen = 7
)


type UpdateUserParams struct {
    FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
}


func (p UpdateUserParams) ToBSON() bson.M {
    m := bson.M{}
    if len(p.FirstName) > 0 {
        m["firstName"] = p.FirstName
    }

    if len(p.LastName) > 0 {
        m["lastName"] = p.LastName
    }
    return m
}


type CreateUserParams struct {
    FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"` 
	Password string `json:"password"`
}


func (params CreateUserParams) Validate() map[string]string {
    errors := map[string]string{}
    if len(params.FirstName) < minFirstNameLen {
        errors["firstName"] = fmt.Sprintf("first name length should be at least %d characters", minFirstNameLen)
    }

    if len(params.LastName) < minLastNameLen {
        errors["lastName"] = fmt.Sprintf("last name length should be at least %d characters", minLastNameLen)
    }

    if len(params.Password) < minPasswordLen {
        errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
    } 

    if !IsEmailValid(params.Email) {
        errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
    }

    return errors
}


func IsValidPassword(encpw, pw string) bool {
    return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}


func IsEmailValid(email string) bool {
        var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
        return emailRegex.MatchString(email)
     }


type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName string `bson:"lastName" json:"lastName"`
	Email string `bson:"email" json:"email"` 
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
    IsAdmin bool `bson:"isAdmin" json:"isAdmin"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
    encrpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
    if err != nil {
        return nil, err
    }
    return &User{
        FirstName: params.FirstName,
        LastName:  params.LastName,
        Email:     params.Email,
        EncryptedPassword: string(encrpw),
    }, nil
}
