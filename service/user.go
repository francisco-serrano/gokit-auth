package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const LoginTemplate = "login.gohtml"

type UserService interface {
	HealthCheck() string
	SendTemplateData() (TemplateRender, error)
	Register(user, pass string) (string, error)
}

type userService struct {
	users map[string]string
}

type TemplateRender struct {
	Metadata  TemplateMetadata
	Variables TemplateVariables
}

type TemplateMetadata struct {
	Name string
}

type TemplateVariables struct {
	Name        string
}

func NewUserService() UserService {
	return &userService{users: make(map[string]string)}
}

func (u userService) HealthCheck() string {
	return "ok"
}

func (u userService) SendTemplateData() (TemplateRender, error) {
	return TemplateRender{
		Metadata:  TemplateMetadata{Name: LoginTemplate},
		Variables: TemplateVariables{Name: "USER"},
	}, nil
}

func (u *userService) Register(user, pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error while hashing pass: %w", err)
	}

	u.users[user] = string(hashedPass)

	return "USER SAVED", nil
}
