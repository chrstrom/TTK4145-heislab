package orderTypes

type OrderStateType int

const (
	Received OrderStateType = iota
	Delegate
	Serving
	Completed
)

type Order struct {
	Floor int
	Dir   int
	Cost  int
}

type OrderStamped struct {
	ID      string
	OrderID int
	Order   Order
}

type HallOrder struct {
	OwnerID       string
	ID            int
	DelegatedToID string
	State         OrderStateType
	Floor, Dir    int
	Costs         map[string]int
}
