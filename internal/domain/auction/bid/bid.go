package bid

import "time"

type Bid struct {
	ID          int64
	AuctionID   int64
	GuildID     int64
	Amount      int64
	IsCancelled bool
	CreatedAt   time.Time
}
