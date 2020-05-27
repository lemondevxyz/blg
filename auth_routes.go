package main

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Login struct {
	UserID   string `validate:"required,min=3"`
	Password string `validate:"required,min=8"`
}

var AUTH_NOT_LOGGED = errors.New("Not logged in")
var AUTH_LOGGED = errors.New("Already logged in")
var AUTH_FAILED = errors.New("The provided credenitals are invalid")

func RouteAuthLogin(c *gin.Context) {

	session := sessions.Default(c)
	/* RouteAuthMiddleware replaces this functionality
	log.Println(session.Get("userid"))
	if session.Get("userid") != nil {
		c.AbortWithError(http.StatusForbidden, AUTH_LOGGED)
		return
	}
	*/

	l := &Login{}

	err := c.Bind(l)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = l.Validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	u := VerifyLogin(l.UserID, l.Password)
	if u == nil {
		c.AbortWithError(http.StatusNotFound, AUTH_FAILED)
		return
	}

	session.Set("userid", u.ID)
	session.Save()

}

func RouteAuthLogout(c *gin.Context) {

	session := sessions.Default(c)
	/* RouteAuthMiddleware replaces this functioanlity
	if session.Get("userid") == nil {
		c.AbortWithError(http.StatusForbidden, AUTH_LOGGED)
		return
	}
	*/

	session.Delete("userid")
	session.Save()

}

func RouteAuthMiddleware(logged bool, fn gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		u := RouteAuthReturnUser(c)

		case1 := (u == nil && !logged)
		case2 := (u != nil && logged)

		if case1 || case2 {
			fn(c)
			return
		} else {
			if u != nil && !logged {
				c.AbortWithError(http.StatusForbidden, AUTH_LOGGED)
			} else if u == nil && logged {
				c.AbortWithError(http.StatusForbidden, AUTH_NOT_LOGGED)
			}
		}
	}
}

func RouteAuthReturnUser(c *gin.Context) (u *User) {

	defer func() {
		c.Set("user", u)
	}()

	session := sessions.Default(c)
	uid := session.Get("userid")

	userid, ok := uid.(uint)
	if !ok {
		return
	}

	u = GetUserById(userid)
	if u != nil {
		u.Password = ""
	}

	return
}

func (l *Login) Validate() error {
	return validate.Struct(l)
}
