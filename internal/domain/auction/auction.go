package auction

import (
	"time"
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
	if a.HighestBid == 0 {
		return a.StartingPrice
	}

	// must be at least 5% above current highest
	return a.HighestBid + (a.HighestBid / 20)
}

func (a *Auction) PlaceBid(guildID int64, amount int64, now time.Time) error {
	if a.Status != Active {
		return ErrAuctionNotActive
	}
	if a.SellerID == guildID {
		return ErrSellerCannotBid
	}
	if amount < a.MinNextBid() {
		return ErrBidTooLow
	}

	a.HighestBid = amount
	a.HighestBidderID = &guildID

	// extend from EndTime, not from now
	if a.EndTime.Sub(now) < extendWindow {
		a.EndTime = a.EndTime.Add(extendWindow)
	}

	return nil
}
