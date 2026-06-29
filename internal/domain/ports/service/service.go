package service

import (
	"time"
)

type PriceOracle interface {
	GetBasePrice(itemID int64) (int64, error)
}

type Clock interface {
	Now() time.Time
}

type Services struct {
	Oracle PriceOracle
	Clock  Clock
}
