package user

import (
	"errors"
	"github.com/go-chi/jwtauth"
	config "github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var TokenAuth *jwtauth.JWTAuth

type User struct {
	ID           rune   `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
	PasswordHash string `json:"-"`
}

func init() {
	TokenAuth = jwtauth.New("HS256", []byte(config.GetString("secret")), nil)
}

func (u *User) GenToken() (string, error) {
	_, token, err := TokenAuth.Encode(jwtauth.Claims{
		"id":  u.ID,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 2)
	u.PasswordHash = string(passwordHash)
	return nil
}
