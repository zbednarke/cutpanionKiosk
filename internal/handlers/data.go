package handlers

import (
	"cutpanionKiosk/internal/cache" // or your module name

	"github.com/gofiber/fiber/v2"
)

func GetAggregatedData(c *fiber.Ctx) error {
	data := cache.GetLatest()
	return c.JSON(data)
}
