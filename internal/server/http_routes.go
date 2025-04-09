package server

func (s *HttpServer) RegisterHttpRoutes() {
	api := s.App.Group("/ads")
	// GET /ads
	api.Get("/", s.GetAds)
	// POST /ads/click
	api.Post("/click", s.handleRecordClick)
	// GET /ads/:id/clicks
	api.Get("/:id/clicks", s.handleGetClickCount)
	// GET /ads/:id/analytics
	api.Get("/:id/analytics", s.handleGetClickAnalytics)
}
