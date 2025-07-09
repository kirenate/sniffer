package main

import (
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"sniffer/logger"
	"sniffer/presentation"

	"github.com/gofiber/fiber/v2"
)

var log = logger.MakeLogger()

func main() {
	addr := "0.0.0.0:443"
	app := fiber.New(fiber.Config{Immutable: true})

	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))
	app.Get("/subdomains", presentation.HandleSubdomains)
	app.Get("/", presentation.HandleListSubdomainsPage)

	err := app.ListenTLS(addr, "./.ssl/ssl_cert.pem", "./.ssl/ssl_cert.pem")
	if err != nil {
		log.Panic().Stack().Err(err).Msg("listener failed")
	}
}
