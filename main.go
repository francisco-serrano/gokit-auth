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

	templateHandler := http.NewServer(
		transport.MakeTemplateEndpoint(svc),
		transport.DecodeRequest,
		transport.EncodeResponseTemplate,
	)

	registerHandler := http.NewServer(
		transport.MakeRegisterEndpoint(svc),
		transport.DecodeRegisterRequest,
		transport.EncodeResponseString,
	)

	app := fiber.New()
	app.Get("/health", adaptor.HTTPHandler(userHandler))
	app.Get("/", adaptor.HTTPHandler(templateHandler))
	app.Post("/register", adaptor.HTTPHandler(registerHandler))

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
