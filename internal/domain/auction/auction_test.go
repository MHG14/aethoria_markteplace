package auction_test

import (
	"errors"
	"testing"
	"time"

	"github.com/MHG14/aethoria_marketplace/internal/domain/auction"
)

func baseAuction() auction.Auction {
	now := time.Now()
	return auction.Auction{
		ID:              1,
		ItemID:          1,
		SellerID:        1,
		StartingPrice:   1000,
		HighestBid:      0,
		HighestBidderID: nil,
		EndTime:         now.Add(24 * time.Hour),
		OriginalEndTime: now.Add(24 * time.Hour),
		Status:          auction.Active,
	}
}

func TestPlaceBid_FirstBid_Success(t *testing.T) {
	a := baseAuction()
	now := time.Now()

	err := a.PlaceBid(2, 1000, now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if a.HighestBid != 1000 {
		t.Errorf("expected highest bid 1000, got %d", a.HighestBid)
	}
	if a.HighestBidderID == nil || *a.HighestBidderID != 2 {
		t.Errorf("expected highest bidder 2, got %v", a.HighestBidderID)
	}
}

func TestPlaceBid_SellerCannotBid(t *testing.T) {
	a := baseAuction()
	now := time.Now()

	err := a.PlaceBid(1, 1100, now)

	if !errors.Is(err, auction.ErrSellerCannotBid) {
		t.Errorf("expected ErrSellerCannotBid, got %v", err)
	}
}

func TestPlaceBid_AuctionNotActive(t *testing.T) {
	a := baseAuction()
	a.Status = auction.Finished
	now := time.Now()

	err := a.PlaceBid(2, 1100, now)

	if !errors.Is(err, auction.ErrAuctionNotActive) {
		t.Errorf("expected ErrAuctionNotActive, got %v", err)
	}
}

func TestPlaceBid_BelowMinimum_NoExistingBid(t *testing.T) {
	a := baseAuction()
	now := time.Now()

	// starting price is 1000, first bid must be >= 1000
	err := a.PlaceBid(2, 999, now)

	if !errors.Is(err, auction.ErrBidTooLow) {
		t.Errorf("expected ErrBidTooLow, got %v", err)
	}
}

func TestPlaceBid_BelowMinimum_WithExistingBid(t *testing.T) {
	a := baseAuction()
	now := time.Now()
	_ = a.PlaceBid(2, 1000, now)

	// min next bid = 1000 + 5% = 1050
	err := a.PlaceBid(3, 1049, now)

	if !errors.Is(err, auction.ErrBidTooLow) {
		t.Errorf("expected ErrBidTooLow, got %v", err)
	}
}

func TestPlaceBid_ExactlyAtMinimum_Success(t *testing.T) {
	a := baseAuction()
	now := time.Now()
	_ = a.PlaceBid(2, 1000, now)

	// min next bid = 1000 + 5% = 1050
	err := a.PlaceBid(3, 1050, now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if a.HighestBid != 1050 {
		t.Errorf("expected highest bid 1050, got %d", a.HighestBid)
	}
}

func TestPlaceBid_AutoExtend_WhenBidInLastFiveMinutes(t *testing.T) {
	a := baseAuction()
	now := time.Now()

	// set end time to 3 minutes from now (within extend window)
	a.EndTime = now.Add(3 * time.Minute)

	_ = a.PlaceBid(2, 1000, now)

	expectedEnd := now.Add(3 * time.Minute).Add(5 * time.Minute)
	diff := a.EndTime.Sub(expectedEnd)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expected end time ~%v, got %v", expectedEnd, a.EndTime)
	}
}

func TestPlaceBid_NoExtend_WhenBidNotInLastFiveMinutes(t *testing.T) {
	a := baseAuction()
	now := time.Now()
	originalEnd := a.EndTime

	_ = a.PlaceBid(2, 1000, now)

	if !a.EndTime.Equal(originalEnd) {
		t.Errorf("end time should not change, got %v, want %v", a.EndTime, originalEnd)
	}
}

func TestMinNextBid_NoExistingBid(t *testing.T) {
	a := baseAuction() // StartingPrice = 1000

	// no bids yet — min is exactly starting price
	if a.MinNextBid() != 1000 {
		t.Errorf("expected 1000, got %d", a.MinNextBid())
	}
}

func TestMinNextBid_WithExistingBid(t *testing.T) {
	a := baseAuction()
	a.HighestBid = 2000
	// 2000 + 5% = 2100
	if a.MinNextBid() != 2100 {
		t.Errorf("expected 2100, got %d", a.MinNextBid())
	}
}

func TestPlaceBid_AutoExtend_ExtendsFromEndTime_NotFromNow(t *testing.T) {
	a := baseAuction()
	now := time.Now()

	// end time is 1 minute away
	a.EndTime = now.Add(1 * time.Minute)

	_ = a.PlaceBid(2, 1000, now)

	expectedEnd := now.Add(6 * time.Minute)
	diff := a.EndTime.Sub(expectedEnd)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expected end ~%v (EndTime+5min), got %v", expectedEnd, a.EndTime)
	}
}
