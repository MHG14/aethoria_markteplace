package item

type Status string

const (
	Available Status = "available"
	Listed    Status = "listed"
	Auctioned Status = "auctioned"
)
