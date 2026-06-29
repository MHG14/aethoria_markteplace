package application

import (
	"context"
	"fmt"
	"time"

	"github.com/MHG14/aethoria_marketplace/internal/domain/auction"
	"github.com/MHG14/aethoria_marketplace/internal/domain/item"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/repository"
)

type CreateAuctionRequest struct {
	ItemID        int64 `json:"item_id"`
	SellerID      int64 `json:"seller_id"`
	StartingPrice int64 `json:"starting_price"`
	DurationHours int   `json:"duration_hours"`
}

type CreateAuctionResponse struct {
	Auction auction.Auction `json:"auction"`
}

func (a *App) CreateAuction(ctx context.Context, req CreateAuctionRequest) (CreateAuctionResponse, error) {
	if req.StartingPrice <= 0 {
		return CreateAuctionResponse{}, fmt.Errorf("%w: starting_price must be positive", ErrInvalidInput)
	}

	duration := 24
	if req.DurationHours > 0 {
		duration = req.DurationHours
	}

	i, err := a.repos.Item.GetByID(ctx, req.ItemID)
	if err != nil {
		return CreateAuctionResponse{}, fmt.Errorf("item not found: %w", ErrNotFound)
	}
	if i.OwnerID != req.SellerID {
		return CreateAuctionResponse{}, fmt.Errorf("%w: you do not own this item", ErrForbidden)
	}
	if !i.CanBeListed() {
		if !i.IsLegendary() {
			return CreateAuctionResponse{}, fmt.Errorf("%w: only legendary items can be auctioned", ErrInvalidInput)
		}
		return CreateAuctionResponse{}, ErrItemNotAvailable
	}

	var auc auction.Auction
	err = a.repos.TxManager.WithTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		i, err := repos.Item.GetByIDForUpdate(ctx, req.ItemID)
		if err != nil {
			return fmt.Errorf("item not found: %w", ErrNotFound)
		}
		if i.OwnerID != req.SellerID {
			return fmt.Errorf("%w: you do not own this item", ErrForbidden)
		}
		if !i.CanBeListed() {
			if !i.IsLegendary() {
				return fmt.Errorf("%w: only legendary items can be auctioned", ErrInvalidInput)
			}
			return ErrItemNotAvailable
		}

		if _, err = repos.Auction.GetActiveByItemID(ctx, req.ItemID); err == nil {
			return ErrActiveAuctionExists
		}

		now := a.svc.Clock.Now()
		endTime := now.Add(time.Duration(duration) * time.Hour)

		auc, err = repos.Auction.Create(ctx, auction.Auction{
			ItemID:          req.ItemID,
			SellerID:        req.SellerID,
			StartingPrice:   req.StartingPrice,
			HighestBid:      0,
			HighestBidderID: nil,
			EndTime:         endTime,
			OriginalEndTime: endTime,
			Status:          auction.Active,
		})
		if err != nil {
			return fmt.Errorf("create auction: %w", err)
		}

		if _, err = repos.Item.UpdateStatus(ctx, req.ItemID, item.Auctioned); err != nil {
			return fmt.Errorf("update item status: %w", err)
		}

		return nil
	})
	if err != nil {
		return CreateAuctionResponse{}, err
	}

	return CreateAuctionResponse{Auction: auc}, nil
}
