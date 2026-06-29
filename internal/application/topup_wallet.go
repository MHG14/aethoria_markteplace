package application

import (
	"context"
	"fmt"

	domainerrors "github.com/MHG14/aethoria_marketplace/internal/domain/error"
	"github.com/MHG14/aethoria_marketplace/internal/domain/guild"
)

type TopUpWalletRequest struct {
	GuildID int64 `json:"guild_id"`
	Amount  int64 `json:"amount"`
}

type TopUpWalletResponse struct {
	Guild guild.Guild `json:"guild"`
}

func (a *App) TopUpWallet(ctx context.Context, req TopUpWalletRequest) (TopUpWalletResponse, error) {
	if req.Amount <= 0 {
		return TopUpWalletResponse{}, fmt.Errorf("%w: amount must be positive", domainerrors.ErrInvalidInput)
	}

	if _, err := a.repos.Guild.GetByID(ctx, req.GuildID); err != nil {
		return TopUpWalletResponse{}, fmt.Errorf("guild not found: %w", domainerrors.ErrNotFound)
	}

	g, err := a.repos.Guild.TopUp(ctx, req.GuildID, req.Amount)
	if err != nil {
		return TopUpWalletResponse{}, fmt.Errorf("top up wallet: %w", err)
	}

	return TopUpWalletResponse{Guild: g}, nil
}
