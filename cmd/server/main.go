package main

import (
	"log"
	"time"

	"cutpanionKiosk/internal/handlers"
	"cutpanionKiosk/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Routes
	app.Get("/api/data", handlers.GetAggregatedData)

	// Static files
	app.Static("/", "./web")

	// Kick off periodic polling
	go func() {
		for {
			services.SyncAll() // pull from calendar + sheets
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Fatal(app.Listen(":8080"))
}
