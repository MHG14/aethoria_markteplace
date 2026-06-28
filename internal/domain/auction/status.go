package auction

type Status string

const (
	Active    Status = "active"
	Finished  Status = "finished"
	Cancelled Status = "cancelled"
)
