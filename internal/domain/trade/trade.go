package trade

import "time"

type Type string

const (
	Auction Type = "auction"
	Listing Type = "listing"
)

type Trade struct {
	ID        int64
	ItemID    int64
	SellerID  int64
	BuyerID   int64
	Price     int64
	Type      Type
	CreatedAt time.Time
}
