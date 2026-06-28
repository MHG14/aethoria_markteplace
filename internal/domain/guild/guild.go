package guild

type Guild struct {
	ID            int64
	Name          string
	TotalMoney    int64
	ReservedMoney int64
	DailyLimit    int64
	DailySpent    int64
}

func (g *Guild) AvailableBalance() int64 { return g.TotalMoney - g.ReservedMoney }
func (g *Guild) CanAfford(amount int64) bool {
	return g.AvailableBalance() >= amount && (g.DailySpent+amount) <= g.DailyLimit
}
