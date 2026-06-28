package auction

import (
	"time"

	"github.com/MHG14/aethoria_markteplace/pkg"
)

const extendWindow = 5 * time.Minute

type Auction struct {
	ID              int64
	ItemID          int64
	SellerID        int64
	StartingPrice   int64
	HighestBid      int64
	HighestBidderID *int64
	EndTime         time.Time
	OriginalEndTime time.Time
	Status          Status
	CreatedAt       time.Time
}

func (a *Auction) MinNextBid() int64 {
	base := a.StartingPrice
	if a.HighestBid > 0 {
		base = a.HighestBid
	}
	return base + (base / 20) // +5%
}

func (a *Auction) PlaceBid(guildID int64, amount int64, now time.Time) error {
	if a.Status != Active {
		return pkg.ErrAuctionNotActive
	}
	if a.SellerID == guildID {
		return pkg.ErrSellerCannotBid
	}
	if amount < a.MinNextBid() {
		return pkg.ErrBidTooLow
	}
	a.HighestBid = amount
	a.HighestBidderID = &guildID
	if a.EndTime.Sub(now) < extendWindow {
		a.EndTime = now.Add(extendWindow)
	}
	return nil
}
