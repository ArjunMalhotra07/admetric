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
