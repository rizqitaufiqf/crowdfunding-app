package user

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
}

type service struct {
	// dependencies to save user to database (need User Repository)
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.ID = uuid.New()
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

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	usr, err := s.repository.FindByEmail(email)
	if err != nil {
		return usr, err
	}
	if usr.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return usr, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(password))
	if err != nil {
		return usr, errors.New("invalid email or password")
	}
	return usr, nil
}
