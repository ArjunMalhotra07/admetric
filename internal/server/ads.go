package server

import "github.com/gofiber/fiber/v2"

func (s *HttpServer) GetAds(c *fiber.Ctx) error {
	ads, err := s.AdService.GetAllAds()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch ads",
		})
	}
	return c.JSON(ads)
}

func (s *HttpServer) GetAnalytics(c *fiber.Ctx) error {
	timeframe := c.Query("timeframe", "1h") // Default 1 hour
	adID := c.Query("ad_id", "")            // Optional ad ID filter

	analytics, err := s.MetricsService.GetAnalytics(adID, timeframe)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch analytics",
		})
	}

	return c.JSON(analytics)
}
