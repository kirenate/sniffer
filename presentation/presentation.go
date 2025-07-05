package presentation

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"sniffer/logger"
	"sniffer/sublister"
)

var log = logger.MakeLogger()

func HandleSubdomains(c *fiber.Ctx) error {
	domain := c.Query("domain")
	subdomains := sublister.Sublister(domain)
	if len(subdomains) == 0 {
		return nil
	}

	return c.JSON(subdomains)
}

func HandleListSubdomainsPage(c *fiber.Ctx) error {
	file, err := os.Open("frontend/listSubdomains.html")
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to open listSubdomains.html")
		return err
	}
	html := make([]byte, 10000)
	_, err = file.Read(html)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to read listSubdomains.html")
		return err
	}
	return c.Type("html").Send(html)
}
