package main

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Firstname   string
	Lastname    string
	Email       string `validate:"required,email" gorm:"unique"`
	Username    string `validate:"required,min=3" gorm:"unique"`
	Description string `validate:"max=280"`
	Password    string `validate:"required,min=8"`
}

const BCRYPT_COST = 14

// Creates a new user
func NewUser(u *User) error {
	if u == nil {
		return errors.New("User pointer is nil")
	}

	newuser := &User{}
	*newuser = *u

	err := newuser.Validate()
	if err != nil {
		return err
	}

	err = newuser.HashPassword()
	if err != nil {
		return err
	}

	return db.Create(newuser).Error
}

// Returns a user by their id
func GetUserById(id uint) *User {

	u := &User{}

	db.Where("id = ?", id).First(u)
	if *u == (User{}) {
		return nil
	}

	return u
}

// Returns a user by their username
func GetUserByUsername(username string) *User {
	u := &User{}

	db.Where("username = ?", username).First(u)
	if *u == (User{}) {
		return nil
	}
	u.Password = ""

	return u
}

// Returns a user by their email/username and password
func VerifyLogin(userid, password string) *User {
	u := &User{}

	db.Where("email = ?", userid).Or("username = ?", userid).First(u)

	if u == nil || *u == (User{}) {
		return nil
	}

	if !u.VerifyPassword(password) {
		return nil
	}

	return u
}

// Updates the user information
func (u *User) Update() error {

	if err := u.Validate(); err != nil {
		return err
	}

	oldu := GetUserById(u.ID)
	if oldu == nil {
		return errors.New("No user exists with id")
	}

	if u.Password != oldu.Password {
		err := u.HashPassword()
		if err != nil {
			return err
		}
	}

	return db.Save(u).Error
}

// Deletes the user
func (u *User) Delete() error {
	return db.Unscoped().Delete(u).Error
}

// Hashes the user's password
func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), BCRYPT_COST)
	if err == nil {
		u.Password = string(bytes)
	}

	return err
}

// Compares UserHash and provided password
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Validates the structure of the user
func (u *User) Validate() error {
	return validate.Struct(u)
}
