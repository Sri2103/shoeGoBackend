package userRepository

import (
	"github.com/hashicorp/go-hclog"
	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
	postgresModels "github.com/sri2103/shoeMart/internal/db/postgresSchemas"
	"gorm.io/gorm"
)

type UserRepoImpl struct {
	logger hclog.Logger
	db     *gorm.DB
}

func NewUserRepo(db *gorm.DB, l hclog.Logger) *UserRepoImpl {
	return &UserRepoImpl{
		logger: l,
		db:     db,
	}
}

func (u *UserRepoImpl) AddUser(user *userModel.User) (*userModel.User, error) {
	var dbUser = new(postgresModels.User)
	err := dbUser.FromEntity(user)
	if err != nil {
		u.logger.Error("error while converting to postgres model")
		return nil, err
	}
	err = u.db.Create(dbUser).Error
	if err != nil {
		u.logger.Error("error in creating the user")
		return nil, err
	}

	return dbUser.ToEntity(), nil
}

func (u *UserRepoImpl) FindUser(email string) (*userModel.User, error) {
	var dbUser = new(postgresModels.User)
	err := u.db.Where(&postgresModels.User{Email: email}).First(dbUser).Error
	if err != nil {
		u.logger.Error("no such user found")
		return nil, err
	}
	return dbUser.ToEntity(), nil
}

func (u *UserRepoImpl) GetUserById(userId string) (*userModel.User, error) {
	var dbUser = new(postgresModels.User)
	err := u.db.Find(dbUser, "id=? ", userId).Error
	if err != nil {
		u.logger.Info("No Such User Found")
		return nil, err
	}
	return dbUser.ToEntity(), nil
}
