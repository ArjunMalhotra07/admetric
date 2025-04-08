package server

import (
	"time"

	"github.com/ArjunMalhotra/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s *HttpServer) handleRecordClick(c *fiber.Ctx) error {
	//! Parse
	var click model.Click
	if err := c.BodyParser(&click); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}
	//! Quick validation
	if click.AdID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ad ID is required",
		})
	}
	if click.PlaybackTime <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ad Playback time must be greater than zero",
		})
	}
	click.ID = uuid.New().String()
	click.IP = c.IP()
	click.Timestamp = time.Now()
	//! Async processing - don't wait for this to complete
	go func() {
		if err := s.ClickService.RecordClick(click); err != nil {
			s.Log.Logger.Errorf("Failed to record click", "error", err)
		}
	}()
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Click recorded",
	})
}
