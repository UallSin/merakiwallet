package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"meraki/finalwallet/internal/model"
	"meraki/finalwallet/internal/repo"
	"meraki/finalwallet/internal/utils"
)

type UserService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(email string, password string) error {

	if len(password) < utils.MinPasswordLen {
		return fmt.Errorf("minimum password length: %v", utils.MinPasswordLen)
	}
	if !utils.ValidEmail(email) {
		return fmt.Errorf("wrong email format")
	}
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		logrus.Errorf("Failed to check email uniqe :%s", err.Error())
		return fmt.Errorf("internal server error")
	}
	if existingUser != nil {
		return fmt.Errorf("email already exists")
	}
	newUser := &model.User{
		Email:    email,
		Password: password,
	}
	if err := s.userRepo.CreateUser(newUser); err != nil {
		logrus.Errorf("failed to create new user: %s", err.Error())
		return fmt.Errorf("Internal server error ")
	}
	return nil
}
