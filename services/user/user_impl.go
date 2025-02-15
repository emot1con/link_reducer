package user

import (
	"go_link_reducer/types"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepositoryImpl(DB *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: DB,
	}
}

func (u *UserRepositoryImpl) Register(userPayload types.RegisterPayload) (string, error) {
	user := types.User{
		Name:     userPayload.Name,
		Email:    userPayload.Email,
		Password: userPayload.Password,
	}

	if err := u.DB.Create(&user).Error; err != nil {
		return "", err
	}
	return "User registered successfully", nil
}

func (u *UserRepositoryImpl) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return types.User{}, err
	}
	return user, nil
}

func (u *UserRepositoryImpl) GetUserByID(ID int) (types.User, error) {
	var user types.User
	if err := u.DB.First(&user, ID).Error; err != nil {
		return types.User{}, err
	}
	return user, nil
}
