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
	Role       string   `json:"role"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Mobile     string   `json:"mobile"`
	ImageUrl   string   `json:"image_url"`
	Modules    []AllowedModule `json:"modules"`
}

type AllowedModule struct {
	ModuleName     string `json:"module_name"`
	AllowedActions string `json:"allowed_actions"`
}

func (result *User) PrepareGive() {
	result.Password = ""
}

type UsersEdge Edge[User]

type UsersConnection struct {
	PageInfo 	*PageInfo    	`json:"pageInfo"`
	Edges    	[]*UsersEdge 	`json:"edges"`
}

// node
// returns decoded curosr string
func (s User) GetCursor() string {
	return s.CreatedAt.String()
}


func Login(ctx context.Context, username string, password string) (*LoginInfo, error) {

	db := config.GetDB()
	var err error
	var result LoginInfo

	user := User{}

	err = db.WithContext(ctx).Model(User{}).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &result, errors.New("invalid username or password")
	}
	err = utils.ComparePassword(user.Password, password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &result, errors.New("invalid username or password")
	}

	isActive := *user.IsActive
	if !isActive {
		return &result, errors.New("user is disabled")
	}
	token, err := utils.JwtGenerate(user.ID)
	result.Token = token
	result.UserId = user.ID
	result.Name = user.Name
	result.Username = user.Username
	result.Email = user.Email
	result.Phone = user.Phone
	result.ImageUrl = user.ImageUrl

	if user.RoleId == 0 {
		return nil, errors.New("please assign role")
	} else {
		var userRole Role
		if err := db.WithContext(ctx).Model(&Role{}).
			Preload("RoleModules").Preload("RoleModules.Module").
			Where("id = ?", user.RoleId).First(&userRole, user.RoleId).Error; err != nil {
			return nil, err
		}
		result.Role = userRole.Name
		var allowedModules []AllowedModule
		for _, rm := range userRole.RoleModules {
			allowedModules = append(allowedModules, AllowedModule{
				ModuleName:     rm.Module.Name,
				AllowedActions: rm.AllowedActions,
			})
		}
		result.Modules = allowedModules
	}

	if err != nil {
		return &result, err
	}

	return &result, nil
}

// destroy current session
func Logout(ctx context.Context) (bool, error) {
	token, ok := utils.GetTokenFromContext(ctx)
	if !ok || token == "" {
		return false, errors.New("token is required")
	}

	// Invalidate the token by storing it in Redis with an expiration
	expiration := time.Hour // Match this with your JWT token's expiration time
	err := config.SetRedisValue(token, "invalid", expiration)
	if err != nil {
		return false, err
	}

	return true, nil
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


func PaginateUser(ctx context.Context, limit *int, after *string,
	name *string, phone *string, mobile *string, email *string, isActive *bool) (*UsersConnection, error) {
	

	db := config.GetDB()

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

	edges, pageInfo, err := FetchPageCompositeCursor[User](db, *limit, after, "created_at", ">")
	if err != nil {
		return nil, err
	}

	var usersConnection UsersConnection

	usersConnection.PageInfo = pageInfo
	
	for _, edge := range edges {
		userEdge := UsersEdge(edge)
		usersConnection.Edges = append(usersConnection.Edges, &userEdge)
	}
	return &usersConnection, err
}