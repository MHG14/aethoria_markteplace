package application

import (
	"context"
	"fmt"

	"github.com/MHG14/aethoria_marketplace/internal/domain/auction/bid"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/repository"
	"github.com/MHG14/aethoria_marketplace/internal/domain/wallet"
)

type PlaceBidRequest struct {
	AuctionID int64 `json:"auction_id"`
	GuildID   int64 `json:"guild_id"`
	Amount    int64 `json:"amount"`
}

type PlaceBidResponse struct {
	Bid bid.Bid `json:"bid"`
}

func (a *App) PlaceBid(ctx context.Context, req PlaceBidRequest) (PlaceBidResponse, error) {
	auc, err := a.repos.Auction.GetByID(ctx, req.AuctionID)
	if err != nil {
		return PlaceBidResponse{}, fmt.Errorf("auction not found: %w", ErrNotFound)
	}

	now := a.svc.Clock.Now()
	if err = auc.PlaceBid(req.GuildID, req.Amount, now); err != nil {
		return PlaceBidResponse{}, err
	}

	g, err := a.repos.Guild.GetByID(ctx, req.GuildID)
	if err != nil {
		return PlaceBidResponse{}, fmt.Errorf("guild not found: %w", ErrNotFound)
	}
	if !g.CanAfford(req.Amount) {
		if g.AvailableBalance() < req.Amount {
			return PlaceBidResponse{}, ErrInsufficientFunds
		}
		return PlaceBidResponse{}, ErrDailyLimitExceeded
	}

	var b bid.Bid
	err = a.repos.TxManager.WithTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		auc, err := repos.Auction.GetByIDForUpdate(ctx, req.AuctionID)
		if err != nil {
			return fmt.Errorf("auction not found: %w", ErrNotFound)
		}

		now := a.svc.Clock.Now()
		prevHighestBid := auc.HighestBid
		prevHighestBidderID := auc.HighestBidderID

		if err = auc.PlaceBid(req.GuildID, req.Amount, now); err != nil {
			return err
		}

		g, err := repos.Guild.GetByIDForUpdate(ctx, req.GuildID)
		if err != nil {
			return fmt.Errorf("guild not found: %w", ErrNotFound)
		}
		if !g.CanAfford(req.Amount) {
			if g.AvailableBalance() < req.Amount {
				return ErrInsufficientFunds
			}
			return ErrDailyLimitExceeded
		}

		// release previous top bidder's reservation
		if prevHighestBidderID != nil && *prevHighestBidderID != req.GuildID {
			prevBidder, err := repos.Guild.GetByIDForUpdate(ctx, *prevHighestBidderID)
			if err != nil {
				return fmt.Errorf("prev bidder not found: %w", err)
			}
			if _, err = repos.Guild.UpdateWallet(ctx,
				prevBidder.ID,
				prevBidder.TotalMoney,
				prevBidder.ReservedMoney-prevHighestBid,
				prevBidder.DailySpent,
			); err != nil {
				return fmt.Errorf("release prev bidder funds: %w", err)
			}
			if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
				GuildID: prevBidder.ID,
				Type:    wallet.TxRelease,
				Amount:  prevHighestBid,
				RefType: wallet.RefAuction,
				RefID:   wallet.RefID(auc.ID),
			}); err != nil {
				return fmt.Errorf("record release transaction: %w", err)
			}
		}

		// reserve funds for new bidder
		if _, err = repos.Guild.UpdateWallet(ctx,
			g.ID,
			g.TotalMoney,
			g.ReservedMoney+req.Amount,
			g.DailySpent+req.Amount,
		); err != nil {
			return fmt.Errorf("reserve funds: %w", err)
		}
		if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
			GuildID: req.GuildID,
			Type:    wallet.TxReserve,
			Amount:  req.Amount,
			RefType: wallet.RefAuction,
			RefID:   wallet.RefID(auc.ID),
		}); err != nil {
			return fmt.Errorf("record reserve transaction: %w", err)
		}

		// persist updated auction
		if _, err = repos.Auction.UpdateBid(ctx,
			auc.ID,
			auc.HighestBid,
			auc.HighestBidderID,
			auc.EndTime,
		); err != nil {
			return fmt.Errorf("update auction: %w", err)
		}

		b, err = repos.Bid.Create(ctx, bid.Bid{
			AuctionID: auc.ID,
			GuildID:   req.GuildID,
			Amount:    req.Amount,
		})
		if err != nil {
			return fmt.Errorf("create bid: %w", err)
		}

		return nil
	})
	if err != nil {
		return PlaceBidResponse{}, err
	}

	return PlaceBidResponse{Bid: b}, nil
}
