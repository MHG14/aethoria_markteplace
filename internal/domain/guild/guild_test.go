package guild_test

import (
	"testing"

	"github.com/MHG14/aethoria_marketplace/internal/domain/guild"
)

func baseGuild() guild.Guild {
	return guild.Guild{
		ID:            1,
		Name:          "Test Guild",
		TotalMoney:    1000,
		ReservedMoney: 200,
		DailyLimit:    500,
		DailySpent:    100,
	}
}

func TestAvailableBalance(t *testing.T) {
	g := baseGuild()
	// 1000 - 200 = 800
	if g.AvailableBalance() != 800 {
		t.Errorf("expected 800, got %d", g.AvailableBalance())
	}
}

func TestCanAfford_Success(t *testing.T) {
	g := baseGuild()
	// available=800, daily remaining=400, amount=300 — should pass both
	if !g.CanAfford(300) {
		t.Error("expected CanAfford to return true")
	}
}

func TestCanAfford_InsufficientBalance(t *testing.T) {
	g := baseGuild()
	// available=800 but amount=900 — fails balance check
	if g.CanAfford(900) {
		t.Error("expected CanAfford to return false due to insufficient balance")
	}
}

func TestCanAfford_DailyLimitExceeded(t *testing.T) {
	g := baseGuild()
	// available=800, but daily remaining = 500-100 = 400, amount=450
	if g.CanAfford(450) {
		t.Error("expected CanAfford to return false due to daily limit")
	}
}

func TestCanAfford_ExactlyAtLimit(t *testing.T) {
	g := baseGuild()
	// daily remaining = 400, amount = 400 — should pass
	if !g.CanAfford(400) {
		t.Error("expected CanAfford to return true at exact daily limit")
	}
}

func TestCanAfford_ZeroAmount(t *testing.T) {
	g := baseGuild()
	if !g.CanAfford(0) {
		t.Error("expected CanAfford to return true for zero amount")
	}
}

func TestCanAfford_FullyReserved(t *testing.T) {
	g := baseGuild()
	g.ReservedMoney = g.TotalMoney // available = 0
	if g.CanAfford(1) {
		t.Error("expected CanAfford to return false when fully reserved")
	}
}
