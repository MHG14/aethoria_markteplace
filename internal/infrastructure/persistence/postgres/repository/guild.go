package repository

import (
	"context"

	"github.com/MHG14/aethoria_marketplace/internal/domain/guild"
	db "github.com/MHG14/aethoria_marketplace/internal/infrastructure/persistence/postgres/sqlc"
)

func (r *guildRepo) Create(ctx context.Context, g guild.Guild) (guild.Guild, error) {
	row, err := r.q.CreateGuild(ctx, &db.CreateGuildParams{
		Name:       g.Name,
		TotalMoney: g.TotalMoney,
		DailyLimit: g.DailyLimit,
	})
	if err != nil {
		return guild.Guild{}, err
	}
	return toGuild(row), nil
}

func (r *guildRepo) TopUp(ctx context.Context, id int64, amount int64) (guild.Guild, error) {
	row, err := r.q.TopUpGuildWallet(ctx, &db.TopUpGuildWalletParams{
		ID:     id,
		Amount: amount,
	})
	if err != nil {
		return guild.Guild{}, err
	}
	return toGuild(row), nil
}

func (r *guildRepo) GetByID(ctx context.Context, id int64) (guild.Guild, error) {
	row, err := r.q.GetGuild(ctx, id)
	if err != nil {
		return guild.Guild{}, err
	}
	return toGuild(row), nil
}

func (r *guildRepo) GetByIDForUpdate(ctx context.Context, id int64) (guild.Guild, error) {
	row, err := r.q.GetGuildForUpdate(ctx, id)
	if err != nil {
		return guild.Guild{}, err
	}
	return toGuild(row), nil
}

func (r *guildRepo) UpdateWallet(ctx context.Context, id, total, reserved, dailySpent int64) (guild.Guild, error) {
	row, err := r.q.UpdateGuildWallet(ctx, &db.UpdateGuildWalletParams{
		ID:            id,
		TotalMoney:    total,
		ReservedMoney: reserved,
		DailySpent:    dailySpent,
	})
	if err != nil {
		return guild.Guild{}, err
	}
	return toGuild(row), nil
}

func toGuild(row *db.Guild) guild.Guild {
	return guild.Guild{
		ID:            row.ID,
		Name:          row.Name,
		TotalMoney:    row.TotalMoney,
		ReservedMoney: row.ReservedMoney,
		DailyLimit:    row.DailyLimit,
		DailySpent:    row.DailySpent,
	}
}
