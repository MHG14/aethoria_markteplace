package wallet

import "time"

type RefType string
type RefID int64

const (
	RefAuction RefType = "auction"
	RefListing RefType = "listing"
	RefTrade   RefType = "trade"
)

type WalletTransaction struct {
	ID        int64
	GuildID   int64
	Type      TxType
	Amount    int64
	RefType   RefType
	RefID     RefID
	CreatedAt time.Time
}
