package service

type UserService interface {
	HealthCheck() string
	SendHTML() (string, error)
	Register(user, pass string) error
}

type userService struct {
}

func NewUserService() UserService {
	return userService{}
}

func (u userService) HealthCheck() string {
	return "ok"
}

func (u userService) SendHTML() (string, error) {
	panic("implement me")
}

func (u userService) Register(user, pass string) error {
	panic("implement me")
}
