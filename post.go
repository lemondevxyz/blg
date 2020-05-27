package main

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	UserID      uint   // used to find the original poster
	Title       string `validate:"required" gorm:"unique"`
	Description string `validate:"required"`
	Content     string `validate:"required"`
	Public      bool
	ImgPath     string // not required
}

// Creates a new post
func NewPost(p *Post) (err error) {

	err = p.Validate()
	if err != nil {
		return
	}

	return db.Create(p).Error
}

// Returns a post using the title as key
func GetPost(title string) *Post {

	p := &Post{}

	db.Where("title = ?", title).First(p)
	if *p == (Post{}) {
		return nil
	}

	return p
}

func GetPostById(id uint) *Post {
	p := &Post{}

	db.Where("id = ?", id).First(p)
	if *p == (Post{}) {
		return nil
	}

	return p
}

// Returns a number of posts, ordered by date of creation
func GetPosts(number int) (ps []*Post) {

	order := "created_at DESC"

	if number > 0 {
		db.Order(order).Limit(number).Find(&ps)
	} else {
		db.Order(order).Find(&ps)
	}

	return
}

func GetPostsByUserId(uid uint) (ps []*Post) {
	order := "created_at DESC"

	db.Order(order).Where("user_id = ?", uid).Find(&ps)

	return
}

func (p *Post) GetOP() *User {
	u := &User{}

	db.Model(p).Related(&u, "UserID")

	return u
}

// Deletes a post
func (p *Post) Delete() error {
	return db.Unscoped().Delete(p).Error
}

// Updates a post
func (p *Post) Update() error {

	if err := p.Validate(); err != nil {
		return err
	}

	return db.Save(p).Error
}

// Validates the post's structure
func (p *Post) Validate() error {

	if p.GetOP() == nil {
		return errors.New("Original poster doesn't exist")
	}

	return validate.Struct(p)
}
