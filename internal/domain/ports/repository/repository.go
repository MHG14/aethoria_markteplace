package repository

import (
	"context"
	"time"

	"github.com/MHG14/aethoria_marketplace/internal/domain/auction"
	"github.com/MHG14/aethoria_marketplace/internal/domain/auction/bid"
	"github.com/MHG14/aethoria_marketplace/internal/domain/guild"
	"github.com/MHG14/aethoria_marketplace/internal/domain/item"
	"github.com/MHG14/aethoria_marketplace/internal/domain/listing"
	"github.com/MHG14/aethoria_marketplace/internal/domain/trade"
	"github.com/MHG14/aethoria_marketplace/internal/domain/wallet"
)

type ItemRepository interface {
	Create(ctx context.Context, i item.Item) (item.Item, error)
	GetByID(ctx context.Context, id int64) (item.Item, error)
	GetByIDForUpdate(ctx context.Context, id int64) (item.Item, error)
	List(ctx context.Context) ([]item.Item, error)
	ListByOwner(ctx context.Context, ownerID int64) ([]item.Item, error)
	UpdateStatus(ctx context.Context, id int64, status item.Status) (item.Item, error)
	UpdateOwner(ctx context.Context, id int64, ownerID int64) (item.Item, error)
}

type GuildRepository interface {
	Create(ctx context.Context, g guild.Guild) (guild.Guild, error)
	GetByID(ctx context.Context, id int64) (guild.Guild, error)
	GetByIDForUpdate(ctx context.Context, id int64) (guild.Guild, error)
	UpdateWallet(ctx context.Context, id int64, total, reserved, dailySpent int64) (guild.Guild, error)
	TopUp(ctx context.Context, id int64, amount int64) (guild.Guild, error)
}

type ListingRepository interface {
	Create(ctx context.Context, l listing.Listing) (listing.Listing, error)
	GetByID(ctx context.Context, id int64) (listing.Listing, error)
	GetByIDForUpdate(ctx context.Context, id int64) (listing.Listing, error)
	UpdateStatus(ctx context.Context, id int64, status listing.Status, buyerID *int64) (listing.Listing, error)
	ListActive(ctx context.Context) ([]listing.Listing, error)
}

type AuctionRepository interface {
	Create(ctx context.Context, a auction.Auction) (auction.Auction, error)
	GetByID(ctx context.Context, id int64) (auction.Auction, error)
	GetByIDForUpdate(ctx context.Context, id int64) (auction.Auction, error)
	GetActiveByItemID(ctx context.Context, itemID int64) (auction.Auction, error)
	UpdateBid(ctx context.Context, id int64, highestBid int64, highestBidderID *int64, endTime time.Time) (auction.Auction, error)
	UpdateStatus(ctx context.Context, id int64, status auction.Status) (auction.Auction, error)
	ListActive(ctx context.Context) ([]auction.Auction, error)
	ListExpired(ctx context.Context) ([]auction.Auction, error)
}

type BidRepository interface {
	Create(ctx context.Context, b bid.Bid) (bid.Bid, error)
	GetByID(ctx context.Context, id int64) (bid.Bid, error)
	Cancel(ctx context.Context, id int64) (bid.Bid, error)
	ListByAuction(ctx context.Context, auctionID int64) ([]bid.Bid, error)
	ListActiveByGuildAndAuction(ctx context.Context, auctionID, guildID int64) ([]bid.Bid, error)
}

type TradeRepository interface {
	Create(ctx context.Context, t trade.Trade) (trade.Trade, error)
	ListByGuild(ctx context.Context, guildID int64) ([]trade.Trade, error)
}

type WalletRepository interface {
	CreateTransaction(ctx context.Context, tx wallet.Transaction) (wallet.Transaction, error)
	ListTransactions(ctx context.Context, guildID int64) ([]wallet.Transaction, error)
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}

type Repositories struct {
	Item      ItemRepository
	Guild     GuildRepository
	Listing   ListingRepository
	Auction   AuctionRepository
	Bid       BidRepository
	Trade     TradeRepository
	Wallet    WalletRepository
	TxManager TxManager
}
