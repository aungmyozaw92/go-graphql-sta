package models

import (
	"context"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Username   string    `gorm:"size:100;not null;unique" json:"username" binding:"required"`
	Name       string    `gorm:"size:100;not null" json:"name" binding:"required"`
	Email      string    `gorm:"size:100;unique" json:"email"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Mobile     string    `gorm:"size:20" json:"mobile"`
	ImageUrl   string    `json:"image_url"`
	Password   string    `gorm:"size:255;not null" json:"password"`
	IsActive   *bool     `gorm:"not null" json:"is_active"`
	RoleId     int       `gorm:"not null;default:0" json:"role_id" binding:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewUser struct {
	Username   string   `json:"username"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Mobile     string   `json:"mobile"`
	ImageUrl   string   `json:"image_url"`
	Password   string   `json:"password"`
	IsActive   *bool    `json:"is_active"`
	RoleId     int      `json:"role_id"`
}

type LoginInfo struct {
	Token      string   `json:"token"`
	UserId	   int      `json:"user_id"`
	Username   string   `json:"username"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Mobile     string   `json:"mobile"`
	ImageUrl   string   `json:"image_url"`
}


func (result *User) PrepareGive() {
	result.Password = ""
}


func Login(ctx context.Context, username string, password string) (*LoginInfo, error) {

	db := config.GetDB()
	var err error
	var result LoginInfo

	u := User{}

	err = db.WithContext(ctx).Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return &result, errors.New("invalid username or password")
	}
	err = utils.ComparePassword(u.Password, password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &result, errors.New("invalid username or password")
	}

	isActive := *u.IsActive
	if !isActive {
		return &result, errors.New("user is disabled")
	}
	token, err := utils.JwtGenerate(u.ID)
	result.Token = token
	result.UserId = u.ID
	result.Name = u.Name
	result.Username = u.Username
	result.Email = u.Email
	result.Phone = u.Phone
	result.ImageUrl = u.ImageUrl
	
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func CreateUser(ctx context.Context, input *NewUser) (*User, error) {

	db := config.GetDB()
	var count int64

	if input.Email != "" && !utils.IsValidEmail(input.Email) {
		return &User{}, errors.New("invalid email address")
	}

	err := db.WithContext(ctx).Model(&User{}).Where("username = ?", input.Username).Or("email = ?", input.Email).Count(&count).Error
	if err != nil {
		return &User{}, err
	}
	if count > 0 {
		return &User{}, errors.New("duplicate username or email")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return &User{}, err
	}

	user := User{
		Username:   html.EscapeString(strings.TrimSpace(input.Username)),
		Name:       input.Name,
		Email:      strings.ToLower(input.Email),
		Phone:      input.Phone,
		Mobile:     input.Mobile,
		ImageUrl:   input.ImageUrl,
		Password:   string(hashedPassword),
		IsActive:   input.IsActive,
		// Role:       input.Role,
		RoleId:     input.RoleId,
	}

	err = db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	user.Password = ""
	return &user, nil
}

func GetUser(ctx context.Context, id int) (*User, error) {

	db := config.GetDB()
	var result User

	err := db.WithContext(ctx).First(&result, id).Error

	if err != nil {
		return &result, errors.New("record not found")
	}

	result.PrepareGive()

	return &result, nil
}