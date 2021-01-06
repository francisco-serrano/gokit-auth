package service

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const MainTemplate = "main.gohtml"

type UserService interface {
	HealthCheck() string
	SendMainTemplateData(token string) (TemplateRender, error)
	Register(user, pass string) (string, error)
	Login(user, pass string) (string, error)
}

type userService struct {
	users    map[string]UserFields
	sessions map[string]string
}

type UserFields struct {
	Username       string
	HashedPassword string
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
	Session      string
	User         string
}

func NewUserService() UserService {
	return &userService{
		users:    make(map[string]UserFields),
		sessions: make(map[string]string),
	}
}

func (u userService) HealthCheck() string {
	return "ok"
}

func (u userService) SendMainTemplateData(token string) (TemplateRender, error) {
	if strings.TrimSpace(token) == "" {
		return TemplateRender{
			Metadata:  TemplateMetadata{Name: MainTemplate},
			Variables: TemplateVariables{},
		}, nil
	}

	sessionID, err := ParseToken(token)
	if err != nil {
		return TemplateRender{
			Metadata:  TemplateMetadata{Name: MainTemplate},
			Variables: TemplateVariables{},
		}, fmt.Errorf("error while parsing token: %w", err)
	}

	user, ok := u.sessions[sessionID]
	if !ok {
		return TemplateRender{
			Metadata:  TemplateMetadata{Name: MainTemplate},
			Variables: TemplateVariables{},
		}, fmt.Errorf("session not registered")
	}

	return TemplateRender{
		Metadata:  TemplateMetadata{Name: MainTemplate},
		Variables: TemplateVariables{Session: token, User: user},
	}, nil
}

func (u *userService) Register(user, pass string) (string, error) {
	if _, ok := u.users[user]; ok {
		return "", fmt.Errorf("user already registered")
	}

	hashedPass, err := u.hashValue(pass)
	if err != nil {
		return "", fmt.Errorf("error while hashing pass: %w", err)
	}

	u.users[user] = UserFields{
		Username:       user,
		HashedPassword: hashedPass,
	}

	return "REGISTER SUCCESSFUL", nil
}

func (u userService) Login(user, pass string) (string, error) {
	userFields, ok := u.users[user]
	if !ok {
		return "", fmt.Errorf("user not registered")
	}

	if err := u.checkPasswordHash(pass, userFields.HashedPassword); err != nil {
		return "", fmt.Errorf("error while checking passwords: %w", err)
	}

	sessionID := uuid.New().String()
	u.sessions[sessionID] = user

	token := CreateToken(sessionID)

	return token, nil
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
