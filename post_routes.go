package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var POST_NOT_OP = errors.New("You need to be original author or admin to continue with this action")
var POST_NONE = errors.New("Post does not exist")

func RoutePostGet(c *gin.Context) {
	p := GetPost(c.Param("title"))
	if p == nil {
		c.AbortWithError(http.StatusBadRequest, POST_NONE)
		return
	}

	c.JSON(200, p)
}

func RoutePostDelete(c *gin.Context) {
	postval, _ := c.Get("post")
	p, ok := postval.(*Post)
	if p == nil || !ok || postval == nil {
		c.AbortWithError(http.StatusBadRequest, POST_NONE)
		return
	}

	err := p.Delete()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

func RoutePostPatch(c *gin.Context) {

	postval, _ := c.Get("post")
	p, ok := postval.(*Post)
	if p == nil || !ok || postval == nil {
		c.AbortWithError(http.StatusBadRequest, POST_NONE)
		return
	}

	bind := &struct {
		Title       string `validate:"required"`
		Description string `validate:"required"`
		Content     string `validate:"required"`
		Public      bool
	}{}

	err := c.Bind(bind)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	case1 := validate.StructExcept(bind, "Description", "Content")
	if case1 == nil {
		p.Title = bind.Title
	}

	case2 := validate.StructExcept(bind, "Title", "Content")
	if case2 == nil {
		p.Description = bind.Description
	}

	case3 := validate.StructExcept(bind, "Title", "Description")
	if case3 == nil {
		p.Content = bind.Content
	}
	p.Public = bind.Public

	err = p.Update()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

// hehe
func RoutePostPost(c *gin.Context) {
	bind := &struct {
		Title       string `binding:"required"`
		Description string `binding:"required"`
		Content     string `binding:"required"`
	}{}

	err := c.Bind(bind)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	u := RouteAuthReturnUser(c)
	if u == nil {
		c.AbortWithError(http.StatusBadRequest, USER_NONE)
		return
	}

	p := &Post{
		Title:       bind.Title,
		Description: bind.Description,
		Content:     bind.Content,
	}
	p.UserID = u.ID

	err = p.Validate()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = NewPost(p)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, p)
}

func RoutePostMiddleware(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := RouteAuthReturnUser(c)
		if u == nil {
			c.AbortWithError(http.StatusForbidden, AUTH_NOT_LOGGED)
			return
		}

		title := c.Param("title")
		p := GetPost(title)
		if p == nil {
			c.AbortWithError(http.StatusNotFound, POST_NONE)
			return
		}

		op := p.GetOP()
		if u == nil || op.ID != u.ID {
			c.AbortWithError(http.StatusForbidden, POST_NOT_OP)
			return
		}

		c.Set("post", p)
		f(c)
	}
}
