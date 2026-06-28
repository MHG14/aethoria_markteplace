package item

type Type string

const (
	Common    Type = "common"
	Rare      Type = "rare"
	Legendary Type = "legendary"
)

type Item struct {
	ID      int64
	Name    string
	Type    Type
	Status  Status
	OwnerID int64
}

func (i *Item) IsLegendary() bool    { return i.Type == Legendary }
func (i *Item) CanBeListed() bool    { return i.Status == Available && i.Type != Legendary }
func (i *Item) CanBeAuctioned() bool { return i.Status == Available && i.Type == Legendary }
