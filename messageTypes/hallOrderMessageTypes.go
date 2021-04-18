package messageTypes

type OrderStateType int

const (
	Received OrderStateType = iota
	Delegate
	Serving
	Completed
)

type OrderStamped struct {
	ID      string
	OrderID int
	Floor int
	Dir int
	Cost int
}

type HallOrder struct {
	OwnerID       string
	ID            int
	DelegatedToID string
	State         OrderStateType
	Floor, Dir    int
	Costs         map[string]int
}
