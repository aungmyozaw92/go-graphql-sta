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
	Email      string    `gorm:"size:100;unique;default:null" json:"email"`
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


func UpdateUser(ctx context.Context, id int, input *NewUser) (*User, error) {

	db := config.GetDB()
	var count int64

	err := db.Model(&User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return &User{}, err
	}
	if count <= 0 {
		return nil, errors.New("record not found")
	}

	if err = db.Model(&User{}).
		Where("username = ? OR email = ?", input.Username, input.Email).
		Not("id = ?", id).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return &User{}, errors.New("duplicate email or username")
	}

	// db action
	var user User
	if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).Model(&user).Updates(map[string]interface{}{
		"Name": input.Name, 
		"Email": input.Email, 
		"Username": input.Username, 
		"Phone": input.Phone, 
		"Mobile": input.Mobile, 
		"IsActive": input.IsActive,
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(ctx context.Context, id int) (*User, error) {

	db := config.GetDB()

	var user User

	err := db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, errors.New("record not found")
	}

	err = db.Delete(&user).Error
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func ChangePassword(ctx context.Context, oldPassword string, newPassword string) (*User, error) {
	userId, ok := utils.GetUserIdFromContext(ctx)
	if !ok || userId == 0 {
		return nil, errors.New("user id is required")
	}

	var user User
	db := config.GetDB()
	if err := db.WithContext(ctx).First(&user, userId).Error; err != nil {
		return nil, err
	}
	// check oldPassword
	if err := utils.ComparePassword(user.Password, oldPassword); err != nil {
		return nil, errors.New("old password is wrong")
	}

	//turn password into hash
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}
	newPassword = string(hashedPassword)

	tx := db.Begin()
	if err := tx.WithContext(ctx).Model(&user).UpdateColumn("password", newPassword).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	return &user, tx.Commit().Error
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

func GetUsers(ctx context.Context,name *string, phone *string, mobile *string, email *string, isActive *bool) ([]*User, error) {

	db := config.GetDB()
	var results []*User

	if name != nil && *name != "" {
		db.Where("name LIKE ?", "%"+*name+"%")
	}
	if phone != nil && *phone != "" {
		db.Where("phone LIKE ?", "%"+*phone+"%")
	}
	if mobile != nil && *mobile != "" {
		db.Where("mobile LIKE ?", "%"+*mobile+"%")
	}
	if email != nil && *email != "" {
		db.Where("email LIKE ?", "%"+*email+"%")
	}
	if isActive != nil {
		db.Where("is_active = ?", isActive)
	}

	if err := db.WithContext(ctx).Find(&results).Error; err != nil {
		return results, errors.New("no user")
	}

	for i, u := range results {
		u.Password = ""
		results[i] = u
	}

	return results, nil
}