package auth

import (
	"golang.org/x/crypto/bcrypt" 
	"errors"
	"go/adv-demo/internal/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService{
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(email, password, name string) (string, error){
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	user := &user.User{
		Email: email,
		Password: string(hashedPassword),
		Name: name,
	}

	_, err = service.UserRepository.Create(user)
	if err != nil{
		return "", err
	}

	return user.Email, nil
}

func (service *AuthService) Login(email, password string) (string, error){
	executedUser, err := service.UserRepository.FindByEmail(email)
	if err != nil{
		return "", errors.New(ErrWrongCredentials)
	} 

	err = bcrypt.CompareHashAndPassword([]byte(executedUser.Password), []byte(password))
	if err != nil{
		return "", errors.New(ErrWrongCredentials)
	}

	return executedUser.Email, nil 
}