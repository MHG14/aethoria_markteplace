package httpserver

func (s *Server) registerRoutes() {
	api := s.app.Group("/api/v1")

	// items
	api.Post("/items", s.handlers.Item.Create)
	api.Get("/items", s.handlers.Item.List)
	api.Get("/items/:id", s.handlers.Item.Get)
	api.Get("/guilds/:id/items", s.handlers.Item.ListByOwner)

	// listings
	api.Post("/listings", s.handlers.Listing.Create)
	api.Get("/listings", s.handlers.Listing.List)
	api.Get("/listings/:id", s.handlers.Listing.Get)
	api.Post("/listings/:id/buy", s.handlers.Listing.Buy)
	api.Delete("/listings/:id", s.handlers.Listing.Cancel)

	// auctions
	api.Post("/auctions", s.handlers.Auction.Create)
	api.Get("/auctions", s.handlers.Auction.List)
	api.Get("/auctions/:id", s.handlers.Auction.Get)
	api.Post("/auctions/:id/bids", s.handlers.Auction.PlaceBid)
	api.Delete("/auctions/:id/bids/:bid_id", s.handlers.Auction.CancelBid)

	// guilds
	api.Post("/guilds", s.handlers.Guild.Create)
	api.Get("/guilds/:id", s.handlers.Guild.Get)
	api.Get("/guilds/:id/wallet", s.handlers.Guild.GetWallet)
	api.Post("/:id/topup", s.handlers.Guild.TopUp)
	api.Get("/guilds/:id/transactions", s.handlers.Guild.GetTransactions)
}
