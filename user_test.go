package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

var userinfo = &User{
	Email:    "test@example.com",
	Username: "test",
	Password: "password",
}

func TestUserValidate(t *testing.T) {
	email := [2]string{
		"dwa", //invalid
		userinfo.Email,
	}

	username := [2]string{
		"aw", //invalid
		userinfo.Username,
	}

	password := [2]string{
		"1234567", //invalid
		userinfo.Password,
	}

	case1 := &User{
		Email:    email[0],
		Username: username[1],
		Password: password[1],
	}

	case2 := &User{
		Email:    email[1],
		Username: username[0],
		Password: password[1],
	}

	case3 := &User{
		Email:    email[1],
		Username: username[1],
		Password: password[0],
	}

	err1 := case1.Validate()
	err2 := case2.Validate()
	err3 := case3.Validate()
	if err1 == nil || err2 == nil || err3 == nil {
		t.Fatalf("One of the user cases was supposed to be invalid, but it's valid. %v\n%v\n%v\n", err1, err2, err3)
	}

	if err := userinfo.Validate(); err != nil {
		t.Fatalf("Failed to validate user: %v", err)
	}
}

func TestNewUser(t *testing.T) {

	err := NewUser(userinfo)
	if err != nil {
		t.Fatalf("An error occured with creating a new user: %v", err)
	}

	twouser := &User{}

	*twouser = *userinfo
	twouser.Username = "test2"

	db.SetLogger(log.New(ioutil.Discard, "\r\n", 0))

	err = NewUser(twouser)
	if err == nil {
		t.Fatal("Email has been taken, but other user still got created.")
	}

	twouser.Username = userinfo.Username
	twouser.Email = "test2@example.com"
	err = NewUser(twouser)
	if err == nil {
		t.Fatal("Username has been taken, but other user still got created.")
	}

	db.SetLogger(gorm.Logger{log.New(os.Stdout, "\r\n", 0)})
}

func TestVerifyLogin(t *testing.T) {
	var u *User
	if u = VerifyLogin(userinfo.Email, userinfo.Password); u == nil {
		t.Fatalf("VerifyLogin failed with email")
	}

	if u = VerifyLogin(userinfo.Username, userinfo.Password); u == nil {
		t.Fatalf("VerifyLogin failed with email")
	}

	userinfo = u
}

func TestUserUpdate(t *testing.T) {
	newpassword := "ayowassupshorty"

	userinfo.Email = "alcool094@gmail.com"
	userinfo.Password = newpassword

	err := userinfo.Update()
	if err != nil {
		t.Fatalf("UserUpdate returned an error: %v", err)
	}

	if VerifyLogin(userinfo.Email, newpassword) == nil {
		t.Fatalf("Password did not update")
	}

}

func TestUserDelete(t *testing.T) {

	id := userinfo.ID

	err := userinfo.Delete()
	if err != nil {
		t.Fatalf("UserDelete returned an error: %v", err)
	}

	if GetUserById(id) != nil {
		t.Fatalf("UserDelete doesn't delete")
	}
}
