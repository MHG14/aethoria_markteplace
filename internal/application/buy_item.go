package application

import (
	"context"
	"fmt"

	"github.com/MHG14/aethoria_marketplace/internal/domain/listing"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/repository"
	"github.com/MHG14/aethoria_marketplace/internal/domain/trade"
	"github.com/MHG14/aethoria_marketplace/internal/domain/wallet"
)

type BuyItemRequest struct {
	ListingID int64 `json:"listing_id"`
	BuyerID   int64 `json:"buyer_id"`
}

type BuyItemResponse struct {
	Trade trade.Trade `json:"trade"`
}

func (a *App) BuyItem(ctx context.Context, req BuyItemRequest) (BuyItemResponse, error) {
	l, err := a.repos.Listing.GetByID(ctx, req.ListingID)
	if err != nil {
		return BuyItemResponse{}, fmt.Errorf("listing not found: %w", ErrNotFound)
	}
	if l.Status != listing.Active {
		return BuyItemResponse{}, ErrItemNotAvailable
	}
	if l.SellerID == req.BuyerID {
		return BuyItemResponse{}, fmt.Errorf("%w: cannot buy your own listing", ErrForbidden)
	}

	buyer, err := a.repos.Guild.GetByID(ctx, req.BuyerID)
	if err != nil {
		return BuyItemResponse{}, fmt.Errorf("buyer not found: %w", ErrNotFound)
	}
	if !buyer.CanAfford(l.Price) {
		if buyer.AvailableBalance() < l.Price {
			return BuyItemResponse{}, ErrInsufficientFunds
		}
		return BuyItemResponse{}, ErrDailyLimitExceeded
	}

	var t trade.Trade
	err = a.repos.TxManager.WithTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		l, err := repos.Listing.GetByIDForUpdate(ctx, req.ListingID)
		if err != nil {
			return fmt.Errorf("listing not found: %w", ErrNotFound)
		}
		if l.Status != listing.Active {
			return ErrItemNotAvailable
		}

		buyer, err := repos.Guild.GetByIDForUpdate(ctx, req.BuyerID)
		if err != nil {
			return fmt.Errorf("buyer not found: %w", ErrNotFound)
		}
		if !buyer.CanAfford(l.Price) {
			if buyer.AvailableBalance() < l.Price {
				return ErrInsufficientFunds
			}
			return ErrDailyLimitExceeded
		}

		seller, err := repos.Guild.GetByIDForUpdate(ctx, l.SellerID)
		if err != nil {
			return fmt.Errorf("seller not found: %w", ErrNotFound)
		}

		if _, err = repos.Guild.UpdateWallet(ctx,
			buyer.ID,
			buyer.TotalMoney-l.Price,
			buyer.ReservedMoney,
			buyer.DailySpent+l.Price,
		); err != nil {
			return fmt.Errorf("deduct buyer wallet: %w", err)
		}

		if _, err = repos.Guild.UpdateWallet(ctx,
			seller.ID,
			seller.TotalMoney+l.Price,
			seller.ReservedMoney,
			seller.DailySpent,
		); err != nil {
			return fmt.Errorf("credit seller wallet: %w", err)
		}

		if _, err = repos.Listing.UpdateStatus(ctx, l.ID, listing.Sold, &req.BuyerID); err != nil {
			return fmt.Errorf("update listing: %w", err)
		}

		if _, err = repos.Item.UpdateOwner(ctx, l.ItemID, req.BuyerID); err != nil {
			return fmt.Errorf("transfer item: %w", err)
		}

		if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
			GuildID: req.BuyerID,
			Type:    wallet.TxDebit,
			Amount:  l.Price,
			RefType: wallet.RefListing,
			RefID:   wallet.RefID(l.ID),
		}); err != nil {
			return fmt.Errorf("record buyer transaction: %w", err)
		}

		if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
			GuildID: l.SellerID,
			Type:    wallet.TxCredit,
			Amount:  l.Price,
			RefType: wallet.RefListing,
			RefID:   wallet.RefID(l.ID),
		}); err != nil {
			return fmt.Errorf("record seller transaction: %w", err)
		}

		t, err = repos.Trade.Create(ctx, trade.Trade{
			ItemID:   l.ItemID,
			SellerID: l.SellerID,
			BuyerID:  req.BuyerID,
			Price:    l.Price,
			Type:     trade.Listing,
		})
		if err != nil {
			return fmt.Errorf("create trade: %w", err)
		}

		return nil
	})
	if err != nil {
		return BuyItemResponse{}, err
	}

	return BuyItemResponse{Trade: t}, nil
}
