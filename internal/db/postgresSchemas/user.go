package postgresModels

import (
	"github.com/google/uuid"
	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primaryKey"`
	Email     string     `gorm:"uniqueIndex"`
	Password  string
	Username  string     `gorm:"uniqueIndex"`
	TokenHash string
}

func (u *User) ToEntity() *userModel.User {
	return &userModel.User{
		ID: u.ID.String(),
		Email: u.Email,
		Password: u.Password,
		Username: u.Username,
		TokenHash: u.Username,
	}
}

func (u *User) FromEntity( user *userModel.User) error {
	if len(user.ID) == 0 {
		u.ID = uuid.New()
	} else {
		id, err := uuid.Parse(user.ID)
		u.ID = id
		if err != nil {
			return err
		}
	}

	u.Email = user.Email
	u.Password= user.Password
	u.Username=user.Username
	u.TokenHash = user.TokenHash
	return nil
}
