package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/fvbock/endless"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type server struct {
	secret string
	port   uint16
	r      *gin.Engine
}

func (s *server) funcMap() template.FuncMap {
	funcmap := sprig.FuncMap()
	myfuncmap := template.FuncMap{
		"urlEq": func(url interface{}, eq string) bool {
			str, ok := url.(string)
			if ok && str == eq {
				return true
			}

			return false
		},

		"getPosts": func(number int) []*Post {
			return GetPosts(number)
		},

		"getPostsByUserId": func(uid uint) []*Post {
			return GetPostsByUserId(uid)
		},

		// prints out html unescaped
		"safeHTML": func(str string) template.HTML {
			return template.HTML(str)
		},

		// gets a template by name and then prints it
		// templates do not change unless you restart
		"getTemplate": func(str string, value interface{}) string {
			buf := new(bytes.Buffer)

			tmpl := template.Must(template.New("").Funcs(funcmap).ParseGlob("templates/**/*"))
			tmpl.ExecuteTemplate(buf, str, value)

			return buf.String()
		},

		"removeQuotes": func(str string) template.JS {
			return template.JS(str)
		},

		"formatDate": func(t time.Time) string {
			return t.Format("Jan, 01, 2006")
		},

		"pathEscape": func(str string) string {
			return url.PathEscape(str)
		},
	}

	for k, v := range myfuncmap {
		funcmap[k] = v
	}

	return funcmap
}

func (s *server) useSession() {
	s.r.Use(sessions.Sessions("sesh", cookie.NewStore([]byte(s.secret))))
}

func (s *server) initRoutes() {

	s.r.SetFuncMap(s.funcMap())
	s.r.LoadHTMLGlob("templates/**/*")

	getParam := func(c *gin.Context) map[string]interface{} {
		var p *Post
		var u *User

		val1, exists := c.Get("Post")
		if exists {
			post, ok := val1.(*Post)
			if ok {
				p = post
			}
		}

		val2, exists := c.Get("User")
		if exists {
			user, ok := val2.(*User)
			if ok {
				u = user
			}
		} else {
			u = RouteAuthReturnUser(c)
		}

		return map[string]interface{}{
			"Post":     p,
			"User":     u,
			"URL":      c.Request.URL.Path,
			"ShareURL": "http://" + c.Request.Host + c.Request.URL.EscapedPath(),
		}
	}

	s.r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", getParam(c))
	})

	// Post-related pages
	s.r.GET("/all", func(c *gin.Context) {
		c.HTML(200, "all.html", getParam(c))
	})

	s.r.GET("/create-post", func(c *gin.Context) {
		c.HTML(200, "create-post.html", getParam(c))
	})

	s.r.GET("/update-post/:title", func(c *gin.Context) {
		c.Set("Post", GetPost(c.Param("title")))
		c.HTML(200, "update-post.html", getParam(c))
	})

	s.r.GET("/view/:title", func(c *gin.Context) {
		title := c.Param("title")
		if len(title) > 0 {
			p := GetPost(title)
			if p != nil {
				c.Set("Post", p)
				c.HTML(200, "single.html", getParam(c))
				return
			}
		}

		c.Status(http.StatusNotFound)
	})

	s.r.GET("/view-by-id/:id", func(c *gin.Context) {
		param := c.Param("id")

		unit, err := strconv.ParseUint(param, 0, 32)
		if err != nil {
			c.Redirect(307, "/")
			return
		}

		id := uint(unit)
		p := GetPostById(id)
		if p == nil {
			c.Redirect(307, "/")
			return
		}

		c.Redirect(307, fmt.Sprintf("/view/%s", p.Title))
	})

	// User-related pages
	s.r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", getParam(c))
	})

	s.r.GET("/logout", func(c *gin.Context) {
		c.HTML(200, "logout.html", getParam(c))
	})

	s.r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", getParam(c))
	})

	s.r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", getParam(c))
	})

	s.r.GET("/profile/:username", func(c *gin.Context) {
		u := GetUserByUsername(c.Param("username"))
		if u != nil {
			param := getParam(c)
			param["SelUser"] = u
			c.HTML(200, "profile.html", param)
		} else {
			c.Redirect(307, "/")
		}
	})

	s.r.GET("/profile-by-id/:id", func(c *gin.Context) {
		param := c.Param("id")

		unit, err := strconv.ParseUint(param, 0, 32)
		if err != nil {
			c.Redirect(307, "/")
			return
		}

		id := uint(unit)
		u := GetUserById(id)
		if u == nil {
			c.Redirect(307, "/")
			return
		}

		c.Redirect(307, fmt.Sprintf("/profile/%s", u.Username))
	})
}

func (s *server) initAPI(router *gin.RouterGroup) {
	// auth_routes.go
	auth := router.Group("/auth")
	{
		auth.POST("/login", RouteAuthMiddleware(false, RouteAuthLogin))
		auth.POST("/logout", RouteAuthMiddleware(true, RouteAuthLogout))

		/*
			auth.GET("/logged", RouteAuthMiddleware(true, func(c *gin.Context) {
				c.String(200, "you should be logged in")
			}))

			auth.GET("/loggedout", RouteAuthMiddleware(false, func(c *gin.Context) {
				c.String(200, "you should be logged out")
			}))
		*/
	}

	// post_routes.go
	post := router.Group("/post")
	{
		post.GET("/:title", RoutePostGet)
		post.DELETE("/:title", RoutePostMiddleware(RoutePostDelete))
		post.PATCH("/:title", RoutePostMiddleware(RoutePostPatch))
		post.POST("/", RoutePostPost)
	}

	// user_routes.go
	user := router.Group("/user")
	{
		user.GET("/", RouteAuthMiddleware(true, RouteUserGet))
		user.DELETE("/", RouteAuthMiddleware(true, RouteUserDelete))
		user.PATCH("/", RouteAuthMiddleware(true, RouteUserPatch))
		user.POST("/", RouteUserPost)
	}

}

func (s *server) Start() error {
	s.useSession()
	s.initAPI(s.r.Group("/api"))
	s.initRoutes()

	err := endless.ListenAndServe(fmt.Sprintf(":%d", s.port), s.r)
	if err != nil {
		return err
	}

	return nil
}
