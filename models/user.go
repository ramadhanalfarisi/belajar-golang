package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRegister struct {
	UserFirstname  string    `json:"userFirstname,omitempty" validate:"required,alpha"`
	UserLastname   string    `json:"userLastname,omitempty" validate:"alpha"`
	UserEmail      string    `json:"userEmail,omitempty" validate:"required,email"`
	UserAddress    string    `json:"userAddress,omitempty"`
	UserPassword   string    `json:"userPassword,omitempty" validate:"required"`
	UserRepassword string    `json:"userRepassword,omitempty" validate:"required,eqfield=UserPassword"`
	UserRole       string    `json:"userRole,omitempty" validate:"required"`
}

type UserLogin struct {
	UserEmail    string `json:"userEmail,omitempty" validate:"required,email"`
	UserPassword string `json:"userPassword,omitempty" validate:"required"`
}

type User struct {
	UserId        uuid.UUID `json:"ID,omitempty" validate:"uuid"`
	UserFirstname string    `json:"userFirstname,omitempty" validate:"required,alpha"`
	UserLastname  string    `json:"userLastname,omitempty" validate:"alpha"`
	UserEmail     string    `json:"userEmail,omitempty" validate:"required,email"`
	UserAddress   string    `json:"userAddress,omitempty"`
	UserPassword  string    `json:"userPassword,omitempty" validate:"required"`
	UserRole      string    `json:"userRole,omitempty" validate:"required"`
}

func (user *User) RegisterUser(db *gorm.DB) error {
	register := db.Create(user)
	if register.Error != nil {
		return register.Error
	}
	return nil
}

func (user *UserLogin) LoginUser(db *gorm.DB) (User, error) {
	var res_user User
	get := db.Where("user_email = ? AND user_password = ?", user.UserEmail, user.UserPassword).Find(&res_user)
	if get.Error != nil {
		return User{}, get.Error
	}
	return res_user, nil
}
