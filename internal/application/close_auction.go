package application

import (
	"context"
	"fmt"

	"github.com/MHG14/aethoria_marketplace/internal/domain/auction"
	"github.com/MHG14/aethoria_marketplace/internal/domain/item"
	"github.com/MHG14/aethoria_marketplace/internal/domain/ports/repository"
	"github.com/MHG14/aethoria_marketplace/internal/domain/trade"
	"github.com/MHG14/aethoria_marketplace/internal/domain/wallet"
)

func (a *App) CloseExpiredAuctions(ctx context.Context) error {
	auctions, err := a.repos.Auction.ListExpired(ctx)
	if err != nil {
		return fmt.Errorf("list expired auctions: %w", err)
	}

	for _, auc := range auctions {
		if err := a.closeAuction(ctx, auc); err != nil {
			fmt.Printf("failed to close auction %d: %v\n", auc.ID, err)
		}
	}

	return nil
}

func (a *App) closeAuction(ctx context.Context, auc auction.Auction) error {
	return a.repos.TxManager.WithTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		auc, err := repos.Auction.GetByIDForUpdate(ctx, auc.ID)
		if err != nil {
			return fmt.Errorf("get auction: %w", err)
		}

		// guard against double-close from concurrent scheduler runs
		if auc.Status != auction.Active {
			return nil
		}

		// no bids — return item to available
		if auc.HighestBidderID == nil {
			if _, err := repos.Auction.UpdateStatus(ctx, auc.ID, auction.Cancelled); err != nil {
				return fmt.Errorf("cancel auction: %w", err)
			}
			if _, err := repos.Item.UpdateStatus(ctx, auc.ItemID, item.Available); err != nil {
				return fmt.Errorf("restore item: %w", err)
			}
			return nil
		}

		winnerID := *auc.HighestBidderID

		winner, err := repos.Guild.GetByIDForUpdate(ctx, winnerID)
		if err != nil {
			return fmt.Errorf("get winner: %w", err)
		}

		// convert winner's reservation to actual debit
		if _, err = repos.Guild.UpdateWallet(ctx,
			winner.ID,
			winner.TotalMoney-auc.HighestBid,
			winner.ReservedMoney-auc.HighestBid,
			winner.DailySpent,
		); err != nil {
			return fmt.Errorf("deduct winner: %w", err)
		}
		if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
			GuildID: winnerID,
			Type:    wallet.TxDebit,
			Amount:  auc.HighestBid,
			RefType: wallet.RefAuction,
			RefID:   wallet.RefID(auc.ID),
		}); err != nil {
			return fmt.Errorf("record winner debit: %w", err)
		}

		seller, err := repos.Guild.GetByIDForUpdate(ctx, auc.SellerID)
		if err != nil {
			return fmt.Errorf("get seller: %w", err)
		}
		if _, err = repos.Guild.UpdateWallet(ctx,
			seller.ID,
			seller.TotalMoney+auc.HighestBid,
			seller.ReservedMoney,
			seller.DailySpent,
		); err != nil {
			return fmt.Errorf("credit seller: %w", err)
		}
		if _, err = repos.Wallet.CreateTransaction(ctx, wallet.Transaction{
			GuildID: auc.SellerID,
			Type:    wallet.TxCredit,
			Amount:  auc.HighestBid,
			RefType: wallet.RefAuction,
			RefID:   wallet.RefID(auc.ID),
		}); err != nil {
			return fmt.Errorf("record seller credit: %w", err)
		}

		if _, err = repos.Item.UpdateOwner(ctx, auc.ItemID, winnerID); err != nil {
			return fmt.Errorf("transfer item: %w", err)
		}

		if _, err = repos.Auction.UpdateStatus(ctx, auc.ID, auction.Finished); err != nil {
			return fmt.Errorf("finish auction: %w", err)
		}

		if _, err = repos.Trade.Create(ctx, trade.Trade{
			ItemID:   auc.ItemID,
			SellerID: auc.SellerID,
			BuyerID:  winnerID,
			Price:    auc.HighestBid,
			Type:     trade.Auction,
		}); err != nil {
			return fmt.Errorf("create trade: %w", err)
		}

		return nil
	})
}
