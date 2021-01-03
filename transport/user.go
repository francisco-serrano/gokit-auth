package transport

import (
	"context"
	"encoding/json"
	"github.com/francisco-serrano/gokit-auth/service"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type healthCheckResponse struct {
	Message string `json:"message"`
}

func MakeHealthEndpoint(svc service.UserService) endpoint.Endpoint {
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return healthCheckResponse{Message: svc.HealthCheck()}, nil
	}
}

func DecodeRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
