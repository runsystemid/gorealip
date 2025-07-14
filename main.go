package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/realclientip/realclientip-go"
)

func realIPMiddleware() fiber.Handler {
	strat, err := realclientip.NewRightmostNonPrivateStrategy("X-Forwarded-For")
	if err != nil {
		log.Fatalf("strategy creation failed: %v", err)
	}

	return func(c *fiber.Ctx) error {
		realIP := strat.ClientIP(c.GetReqHeaders(), c.Context().RemoteAddr().String())
		chain := c.Get("X-Forwarded-For")
		log.Printf("RealIP: %s | ProxyChain: \"%s\" | %s %s %s",
			realIP, chain, c.Method(), c.Path(), c.Protocol())

		log.Printf("context IP: %s\n", c.IPs())
		return c.Next()
	}
}

func main() {
	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Use(realIPMiddleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello\n")
	})

	log.Println("Listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
