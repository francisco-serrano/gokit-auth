package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

const MainTemplate = "main.gohtml"
const LoginTemplate = "login.gohtml"

type UserService interface {
	HealthCheck() string
	SendMainTemplateData() TemplateRender
	Register(user, pass string) (string, error)
	Login(user, pass string) TemplateRender
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
	Name         string
	LoginMessage string
	ErrorMessage error
}

func NewUserService() UserService {
	return &userService{users: make(map[string]string)}
}

func (u userService) HealthCheck() string {
	return "ok"
}

func (u userService) SendMainTemplateData() TemplateRender {
	return TemplateRender{
		Metadata:  TemplateMetadata{Name: MainTemplate},
		Variables: TemplateVariables{},
	}
}

func (u *userService) Register(user, pass string) (string, error) {
	if _, ok := u.users[user]; ok {
		return "", fmt.Errorf("user already registered")
	}

	hashedPass, err := u.hashValue(pass)
	if err != nil {
		return "", fmt.Errorf("error while hashing pass: %w", err)
	}

	u.users[user] = hashedPass

	return "REGISTER SUCCESSFUL", nil
}

func (u userService) Login(user, pass string) TemplateRender {
	storedPass, ok := u.users[user]
	if !ok {
		return TemplateRender{
			Metadata: TemplateMetadata{Name: LoginTemplate},
			Variables: TemplateVariables{
				Name:         user,
				LoginMessage: "LOGIN FAILED",
				ErrorMessage: fmt.Errorf("user not registered"),
			},
		}
	}

	if err := u.checkPasswordHash(pass, storedPass); err != nil {
		return TemplateRender{
			Metadata: TemplateMetadata{Name: LoginTemplate},
			Variables: TemplateVariables{
				Name:         user,
				LoginMessage: "LOGIN FAILED",
				ErrorMessage: fmt.Errorf("erorr while checking passwords: %w", err),
			},
		}
	}

	return TemplateRender{
		Metadata: TemplateMetadata{Name: LoginTemplate},
		Variables: TemplateVariables{
			Name:         user,
			LoginMessage: "LOGIN SUCCESSFUL",
			ErrorMessage: nil,
		},
	}
}

func (u userService) hashValue(v string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (u userService) checkPasswordHash(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
