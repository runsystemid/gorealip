package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/realclientip/realclientip-go"
	"github.com/runsystemid/golog"
)

func getIPFromXForwardedFor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		c.Locals("srcIP", c.Get("X-Forwarded-For"))

		var err error
		if _err := c.Next(); _err != nil {
			err = _err
		}

		logMsg := golog.LogModel{
			Header:       c.Request().Header.Header(),
			Request:      c.Body(),
			Method:       c.Method(),
			HttpStatus:   uint64(c.Response().StatusCode()),
			StatusCode:   strconv.Itoa(c.Response().StatusCode()),
			Response:     string(c.Response().Body()),
			ResponseTime: time.Since(startTime),
			Error:        err,
		}
		golog.TDR(c.Context(), logMsg)

		return err
	}
}

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
		return c.Next()
	}
}

func main() {
	loggerConfig := golog.Config{
		App:             "goreal",
		AppVer:          "1.0.0",
		Env:             "development",
		FileLocation:    "logs/system.log",
		FileTDRLocation: "logs/tdr.log",
		FileMaxSize:     10,
		FileMaxBackup:   10,
		FileMaxAge:      10,
		Stdout:          true,
	}
	golog.Load(loggerConfig)

	app := fiber.New(fiber.Config{
		// ProxyHeader:             fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
	})

	app.Use(getIPFromXForwardedFor())
	app.Use(realIPMiddleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello\n")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK\n")
	})

	log.Println("Listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
