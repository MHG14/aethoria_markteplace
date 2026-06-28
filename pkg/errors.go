package pkg

import "errors"

var (
	ErrAuctionNotActive   = errors.New("auction is not active")
	ErrSellerCannotBid    = errors.New("seller cannot bid on own item")
	ErrBidTooLow          = errors.New("bid must be at least 5% above current highest")
	ErrInsufficientFunds  = errors.New("insufficient available balance")
	ErrDailyLimitExceeded = errors.New("daily purchase limit exceeded")
)
