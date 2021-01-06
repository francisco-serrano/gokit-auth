package main

import (
	"github.com/francisco-serrano/gokit-auth/service"
	"github.com/francisco-serrano/gokit-auth/transport"
	"github.com/go-kit/kit/transport/http"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	svc := service.NewUserService()

	userHandler := http.NewServer(
		transport.MakeHealthEndpoint(svc),
		transport.DecodeRequest,
		transport.EncodeResponseJSON,
	)

	mainHandler := http.NewServer(
		transport.MakeMainEndpoint(svc),
		transport.DecodeRequest,
		transport.SetMainResponse,
	)

	registerHandler := http.NewServer(
		transport.MakeRegisterEndpoint(svc),
		transport.DecodeLoginRegisterRequest,
		transport.EncodeResponseString,
	)

	loginHandler := http.NewServer(
		transport.MakeLoginEndpoint(svc),
		transport.DecodeLoginRegisterRequest,
		transport.SetLoginResponse,
	)

	logoutHandler := http.NewServer(
		transport.MakeLogoutEndpoint(svc),
		transport.DecodeRequest,
		transport.SetLogoutResponse,
	)

	app := fiber.New()
	app.Get("/health", adaptor.HTTPHandler(userHandler))
	app.Get("/", adaptor.HTTPHandler(mainHandler))
	app.Post("/register", adaptor.HTTPHandler(registerHandler))
	app.Post("/login", adaptor.HTTPHandler(loginHandler))
	app.Post("/logout", adaptor.HTTPHandler(logoutHandler))

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
