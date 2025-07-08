package main

import (
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"sniffer/logger"
	"sniffer/presentation"

	//"sniffer/sublister"
	fiber "github.com/gofiber/fiber/v2"
)

var log = logger.MakeLogger()

func main() {
	addr := "localhost:8080"
	app := fiber.New(fiber.Config{Immutable: true})

	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))
	app.Get("/subdomains", presentation.HandleSubdomains)
	app.Get("/", presentation.HandleListSubdomainsPage)

	err := app.ListenTLS(addr, ".ssl/ssl.cert", ".ssl/ssl.key")
	if err != nil {
		log.Panic().Stack().Msg("listener failed")
	}
}
