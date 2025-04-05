package server

func (s *Server) RegisterRoutes() {
	s.Web.RegisterHttpRoutes()
}
