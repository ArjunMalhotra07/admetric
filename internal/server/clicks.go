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

func (s *HttpServer) handleGetClickCount(c *fiber.Ctx) error {
	adID := c.Params("id")
	if adID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ad ID is required",
		})
	}

	// First check if ad exists
	exists, err := s.ClickService.AdExists(adID)
	if err != nil {
		s.Log.Logger.Errorf("Failed to check if ad exists: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ad not found",
		})
	}

	// Get click count from in-memory counter or database
	count, err := s.ClickService.GetClickCount(adID)
	if err != nil {
		s.Log.Logger.Errorf("Failed to get click count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"ad_id":        adID,
		"total_clicks": count,
	})
}
