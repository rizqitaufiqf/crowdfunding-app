package user

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserDTO) (User, error)
	Login(input LoginDTO) (User, error)
	IsEmailAvailable(input CheckEmailDTO) (bool, error)
	SaveAvatar(ID string, imageLocation string) (User, error)
	GetUserByID(ID string) (User, error)
}

type service struct {
	// dependencies to save user to database (need User Repository)
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserDTO) (User, error) {
	user := User{}
	user.ID = uuid.New().String()
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *service) Login(input LoginDTO) (User, error) {
	email := input.Email
	password := input.Password
	usr, err := s.repository.FindByEmail(email)
	if err != nil {
		return usr, err
	}

	if usr.ID == "00000000-0000-0000-0000-000000000000" {
		return usr, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password))
	if err != nil {
		return usr, errors.New("invalid email or password")
	}

	return usr, nil
}

func (s *service) IsEmailAvailable(input CheckEmailDTO) (bool, error) {
	email := input.Email
	usr, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if usr.ID == "00000000-0000-0000-0000-000000000000" {
		return true, nil
	}

	return false, nil
}

func (s *service) SaveAvatar(ID string, imageLocation string) (User, error) {
	usr, err := s.repository.FindByID(ID)
	if err != nil {
		return usr, err
	}

	usr.AvatarFileName = imageLocation
	updateUser, err := s.repository.Update(usr)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (s *service) GetUserByID(ID string) (User, error) {
	usr, err := s.repository.FindByID(ID)
	if err != nil {
		return usr, err
	}

	if usr.ID == "00000000-0000-0000-0000-000000000000" {
		return usr, errors.New("user not found")
	}

	return usr, nil

}
