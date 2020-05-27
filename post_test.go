package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

var posts = []*Post{}

func AppendPosts() {
	for i := 0; i < 5; i++ {

		post := &Post{
			Title:       fmt.Sprintf("post number %d", i),
			Description: fmt.Sprintf("post description %d", i),
			Content:     "yo",
		}

		posts = append(posts, post)
	}
}

func ComparePost(p *Post, v *Post) (b bool) {

	if p != nil && v != nil {
		if p.Title == v.Title && p.Description == v.Description && p.Content == v.Content {
			return true
		}
	}

	return
}

func TestPostValidate(t *testing.T) {
	var val = "string"
	one := &Post{
		Description: val,
		Content:     val,
	}

	if err := one.Validate(); err == nil {
		t.Fatalf("Validation is failing constraint: Title, required.\n%v", one)
	}

	two := &Post{
		Title:   val,
		Content: val,
	}

	if err := two.Validate(); err == nil {
		t.Fatalf("Validation is failing constraint: Description, required.\n%v", two)
	}

	three := &Post{
		Title:       val,
		Description: val,
	}

	if err := three.Validate(); err == nil {
		t.Fatalf("Validation is failing constraint: Content, required.\n%v", two)
	}

	four := &Post{
		Title:       val,
		Content:     val,
		Description: val,
	}

	if err := four.Validate(); err != nil {
		t.Fatalf("Validation is failing constraint, error: %v", err)
	}
}

func TestNewPost(t *testing.T) {

	validate = validator.New()

	CreateDatabase("db.sqlite3")

	AppendPosts()

	for i, v := range posts {

		if err := NewPost(v); err != nil {
			t.Fatalf("An error occured with creating the post, number: %d\nerror: %v", i, err)
		}

		check := Post{}
		time.Sleep(time.Millisecond * 10)

		db.Where(Post{Title: v.Title}).First(&check)

		//t.Logf("post number: %d, post title: %s, post description: %s, post content: %s", i, v.Title, v.Description, v.Content)
		if check == (Post{}) {
			t.Fatalf("NewPost doesn't actually create the post, number: %d\n", i)
		}

		db.SetLogger(log.New(ioutil.Discard, "\r\n", 0))
		if err := NewPost(v); err == nil {
			t.Fatalf("Failing constraint unique for post title, number: %d", i)
		}
		db.SetLogger(gorm.Logger{log.New(os.Stdout, "\r\n", 0)})

	}
}

func TestGetPost(t *testing.T) {
	for _, v := range posts {

		p := GetPost(v.Title)
		if !ComparePost(p, v) {
			t.Fatalf("GetPost doesn't get post by title, value: %v\npost: %v", v, p)
		}

		f := GetPost(v.Description)
		if f != nil {
			t.Fatalf("GetPost gets post by description not title, value: %v\npost: %v", v, f)
		}
	}
}

func TestGetPosts(t *testing.T) {
	allposts := GetPosts(0)
	threeposts := GetPosts(3)

	if len(allposts) != len(posts) {
		t.Logf("GetPosts doesn't isn't the same as defined posts, posts: %d\nallposts: %d", len(posts), len(allposts))
	}

	organizedposts := []*Post{}

	for i := len(allposts); i > 0; i-- {
		organizedposts = append(organizedposts, allposts[i-1])
	}

	for k := range posts {
		p, o := posts[k], organizedposts[k]

		if !ComparePost(p, o) {
			t.Fatalf("GetPosts isn't organized by date, number: %d", k)
		}
	}

	for k, v := range threeposts {
		if !ComparePost(v, allposts[k]) {
			t.Fatalf("GetPosts with number limit isn't organized by date, number: %d", k)
		}
	}
}

func TestPostUpdate(t *testing.T) {

	for i, v := range posts {
		newtitle := "new " + v.Title

		p := GetPost(v.Title)
		if p == nil {
			t.Fatalf("PostUpdate cannot get post from title")
		}

		p.Title = newtitle
		p.Update()

		p = GetPost(newtitle)
		if p == nil {
			t.Fatalf("PostUpdate doesn't update the post, number: %v\npost: %v", i, v)
		}
	}
}

func TestPostDelete(t *testing.T) {
	for _, v := range posts {
		p := GetPost("new " + v.Title)
		if p == nil {
			t.Fatalf("PostDelete cannot get newly updated post")
		}

		p.Delete()

		if GetPost(v.Title) != nil {
			t.Fatalf("PostDelete doesn't delete the post")
		}
	}
}
