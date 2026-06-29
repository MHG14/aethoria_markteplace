package application

import (
	"context"
	"fmt"

	"github.com/MHG14/aethoria_marketplace/internal/domain/auction"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/repository"
	"github.com/MHG14/aethoria_marketplace/internal/domain/wallet"
)

type CancelBidRequest struct {
	AuctionID int64 `json:"auction_id"`
	BidID     int64 `json:"bid_id"`
	GuildID   int64 `json:"guild_id"`
}

func (a *App) CancelBid(ctx context.Context, req CancelBidRequest) error {
	auc, err := a.repos.Auction.GetByID(ctx, req.AuctionID)
	if err != nil {
		return fmt.Errorf("auction not found: %w", ErrNotFound)
	}
	if auc.Status != auction.Active {
		return ErrAuctionNotActive
	}
	if auc.HighestBidderID != nil && *auc.HighestBidderID == req.GuildID {
		return ErrCannotCancelTopBid
	}

	b, err := a.repos.Bid.GetByID(ctx, req.BidID)
	if err != nil {
		return fmt.Errorf("bid not found: %w", ErrNotFound)
	}
	if b.GuildID != req.GuildID {
		return fmt.Errorf("%w: you do not own this bid", ErrForbidden)
	}
	if b.IsCancelled {
		return fmt.Errorf("%w: bid already cancelled", ErrInvalidInput)
	}

	return a.repos.TxManager.WithTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		auc, err := repos.Auction.GetByIDForUpdate(ctx, req.AuctionID)
		if err != nil {
			return fmt.Errorf("auction not found: %w", ErrNotFound)
		}
		if auc.Status != auction.Active {
			return ErrAuctionNotActive
		}
		if auc.HighestBidderID != nil && *auc.HighestBidderID == req.GuildID {
			return ErrCannotCancelTopBid
		}

		if _, err = repos.Bid.Cancel(ctx, req.BidID); err != nil {
			return fmt.Errorf("cancel bid: %w", err)
		}

		g, err := repos.Guild.GetByIDForUpdate(ctx, req.GuildID)
		if err != nil {
			return fmt.Errorf("guild not found: %w", ErrNotFound)
		}
		if _, err = repos.Guild.UpdateWallet(ctx,
			g.ID,
			g.TotalMoney,
			g.ReservedMoney-b.Amount,
			g.DailySpent,
		); err != nil {
			return fmt.Errorf("release funds: %w", err)
		}

		if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
			GuildID: req.GuildID,
			Type:    wallet.TxRelease,
			Amount:  b.Amount,
			RefType: wallet.RefAuction,
			RefID:   wallet.RefID(req.AuctionID),
		}); err != nil {
			return fmt.Errorf("record release transaction: %w", err)
		}

		return nil
	})
}
