package httpserver

func (s *Server) startRouter() {
	s.gin.GET("/current", s.controller.GetActivity)
	s.gin.GET("/current/br", s.controller.GetBrazilActivity)
	s.gin.GET("/history/24h", s.controller.Get24hHistory)
	s.gin.GET("/history/:minutes", s.controller.GetHistoryByMinutes)
	s.gin.GET("/history/br/24h", s.controller.GetBrazil24hHistory)
	s.gin.GET("/history/br/:minutes", s.controller.GetBrazilHistoryByMinutes)
	s.gin.POST("/current", s.controller.saveCurrent)
	// s.gin.GET("/current/geo", s.controller.GetGeoActivity)
	// s.gin.GET("/history/geo/24h", s.controller.GetGeo24hHistory)
	// s.gin.GET("/history/geo/:minutes", s.controller.GetGeoHistoryByMinutes)
}
