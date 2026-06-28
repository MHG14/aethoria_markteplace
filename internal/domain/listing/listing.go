package listing

import "time"

type Listing struct {
	ID        int64
	ItemID    int64
	SellerID  int64
	BuyerID   *int64
	Price     int64
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
