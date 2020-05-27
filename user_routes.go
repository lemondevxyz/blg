package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	USER_NO_ID = errors.New("please provide an id")
	USER_NONE  = errors.New("user doesn't exist")
)

func RouteUserGet(c *gin.Context) {

	userval, _ := c.Get("user")
	u, ok := userval.(*User)
	if u == nil || !ok || userval == nil {
		c.AbortWithError(http.StatusBadRequest, USER_NONE)
		return
	}

	u.Password = ""

	c.JSON(200, u)
}

func RouteUserDelete(c *gin.Context) {

	userval, _ := c.Get("user")
	u, ok := userval.(*User)
	if u == nil || !ok || userval == nil {
		c.AbortWithError(http.StatusBadRequest, USER_NONE)
		return
	}

	u.Delete()
}

func RouteUserPatch(c *gin.Context) {

	userval, _ := c.Get("user")
	u, ok := userval.(*User)
	if u == nil || !ok || userval == nil {
		c.AbortWithError(http.StatusBadRequest, USER_NONE)
		return
	}

	u = GetUserById(u.ID)
	if u == nil {
		c.AbortWithError(http.StatusBadRequest, USER_NONE)
		return
	}

	newuser := &User{}
	err := c.Bind(newuser)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// validating for email
	case1 := validate.StructExcept(newuser, "Username", "Password")
	if case1 == nil {
		u.Email = newuser.Email
	}

	// validating for username
	case2 := validate.StructExcept(newuser, "Email", "Password")
	if case2 == nil {
		u.Username = newuser.Username
	}

	// validating for password
	case3 := validate.StructExcept(newuser, "Email", "Username")
	if case3 == nil {
		u.Password = newuser.Password
	}

	if newuser.Description != u.Description && len(newuser.Description) > 0 {
		u.Description = newuser.Description
	}

	if newuser.Firstname != u.Firstname && len(newuser.Firstname) > 0 {
		u.Firstname = newuser.Firstname
	}

	if newuser.Lastname != u.Lastname && len(newuser.Lastname) > 0 {
		u.Lastname = newuser.Lastname
	}

	err = u.Update()
	if err != nil {
		c.JSON(500, err)
		return
	}
}

func RouteUserPost(c *gin.Context) {

	bind := &struct {
		Email    string `binding:"required"`
		Username string `binding:"required"`
		Password string `binding:"required"`
	}{}

	err := c.Bind(bind)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	u := &User{
		Email:    bind.Email,
		Username: bind.Username,
		Password: bind.Password,
	}

	err = u.Validate()
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	err = NewUser(u)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	user := &User{}
	*user = *u
	user.Password = ""

	c.JSON(200, user)
}
