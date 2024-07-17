package types

import (
    "fmt"
    "regexp"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const (
    bcryptCost = 12
    minFirstNameLen = 2
    minLastNameLen = 2
    minPasswordLen = 7
)

type CreateUserParams struct {
    FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"` 
	Password string `json:"password"`
 
}


func (params CreateUserParams) Validate() []string {
    errors := []string{}
    if len(params.FirstName) < minFirstNameLen {
        errors = append(errors, fmt.Sprintf("first name length should be at least %d characters", minFirstNameLen))
    }

    if len(params.LastName) < minLastNameLen {
        errors = append(errors, fmt.Sprintf("last name length should be at least %d characters", minLastNameLen))
    }

    if len(params.Password) < minPasswordLen {
        errors = append(errors, fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
    } 

    if !IsEmailValid(params.Email) {
        errors = append(errors, fmt.Sprintf("email is invalid"))
    }

    return errors
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
