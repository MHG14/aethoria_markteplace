package handlers

import "github.com/MHG14/aethoria_marketplace/internal/application"

type Handlers struct {
	Item    *ItemHandler
	Listing *ListingHandler
	Auction *AuctionHandler
	Guild   *GuildHandler
}

func New(app *application.App) *Handlers {
	return &Handlers{
		Item:    NewItemHandler(app),
		Listing: NewListingHandler(app),
		Auction: NewAuctionHandler(app),
		Guild:   NewGuildHandler(app),
	}
}
