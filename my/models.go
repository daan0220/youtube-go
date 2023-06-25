package my

import (
	"github.com/jinzhu/gorm"
)

// USER model

type User struct {
	gorm.Model

	Account  string
	Name     string
	Password string
	Message  string
}

// Post model

type Post struct {
	gorm.Model

	Address string
	Message string
	UserId  int
	GroupId int
}

// Group model

type Group struct {
	gorm.Model

	Name    string
	Message string
	UserId  int
}

// Comment model

type Comment struct {
	gorm.Model

	PostId  int
	Message string
	UserId  int
}

// CommentJoin join model

type CommentJoin struct {
	gorm.Model
	Comment
	User
	Post
}
