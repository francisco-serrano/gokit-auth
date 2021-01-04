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

func MakeTemplateEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return svc.SendMainTemplateData(), nil
	}
}

func MakeRegisterEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		userData, ok := request.(loginRegisterRequest)
		if !ok {
			return nil, fmt.Errorf("erorr while casting to register request: %T", request)
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
			return nil, fmt.Errorf("erorr while casting to register request: %T", request)
		}

		return svc.Login(userData.User, userData.Pass), nil
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

func EncodeResponseTemplate(_ context.Context, w http.ResponseWriter, response interface{}) error {
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
