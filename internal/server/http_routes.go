package server

func (s *HttpServer) RegisterHttpRoutes() {
	api := s.App.Group("/ads")
	// GET /ads
	api.Get("/", s.GetAds)
	// POST /ads/click
	api.Post("/click", s.handleRecordClick)
	// GET /ads/analytics
	api.Get("/analytics", s.GetAnalytics)
}
