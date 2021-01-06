package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/francisco-serrano/gokit-auth/service"
	"github.com/go-kit/kit/endpoint"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type healthCheckResponse struct {
	Message string `json:"message"`
}

type loginRegisterRequest struct {
	User string
	Pass string
}

func MakeHealthEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return healthCheckResponse{Message: svc.HealthCheck()}, nil
	}
}

func MakeMainEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		c, ok := request.(*http.Cookie)
		if !ok {
			return nil, fmt.Errorf("could not obtain cookie from request: %T", request)
		}

		render, err := svc.SendMainTemplateData(c.Value)
		if err != nil {
			log.Print(fmt.Errorf("error while obtaining render: %w", err))
		}

		return render, nil
	}
}

func MakeRegisterEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		userData, ok := request.(loginRegisterRequest)
		if !ok {
			return nil, fmt.Errorf("error while casting to register request: %T", request)
		}

		response, err := svc.Register(userData.User, userData.Pass)
		if err != nil {
			return nil, fmt.Errorf("error while registering email: %w", err)
		}

		return response, nil
	}
}

func MakeLoginEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		userData, ok := request.(loginRegisterRequest)
		if !ok {
			log.Print(fmt.Errorf("error while casting to register request: %T", request))

			return "", nil
		}

		token, err := svc.Login(userData.User, userData.Pass)
		if err != nil {
			log.Print(fmt.Errorf("error during login: %w", err))

			return "", nil
		}

		return token, nil
	}
}

func MakeLogoutEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		c, ok := request.(*http.Cookie)
		if !ok {
			return nil, fmt.Errorf("could not obtain cookie from request: %T", request)
		}

		if err := svc.Logout(c.Value); err != nil {
			log.Print(fmt.Errorf("error while logging out: %w", err))
		}

		return nil, nil
	}
}

func DecodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}

	return c, nil
}

func DecodeLoginRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	user := r.FormValue("user")
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("cannot register an empty user")
	}

	pass := r.FormValue("pass")
	if strings.TrimSpace(pass) == "" {
		return nil, fmt.Errorf("cannot register an empty password")
	}

	return loginRegisterRequest{
		User: user,
		Pass: pass,
	}, nil
}

func EncodeResponseJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeResponseString(_ context.Context, w http.ResponseWriter, _ interface{}) error {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func SetMainResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("content-type", "text/html")

	tr, ok := response.(service.TemplateRender)
	if !ok {
		return fmt.Errorf("error while casting template response: %T", response)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("error while geting cwd: %w", err))
	}

	parsedTemplate, err := template.ParseFiles(filepath.Join(cwd, "templates", tr.Metadata.Name))
	if err != nil {
		log.Fatal(fmt.Errorf("error while parsing template: %w", err))
	}

	if err := parsedTemplate.Execute(w, tr.Variables); err != nil {
		return fmt.Errorf("error while executing template: %w", err)
	}

	return nil
}

func SetLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	token, ok := response.(string)
	if !ok {
		return fmt.Errorf("error while casting login response: %T", response)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: token,
	})

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func SetLogoutResponse(_ context.Context, w http.ResponseWriter, _ interface{}) error {
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}
