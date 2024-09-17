package model

type StateType int

const (
	NEW StateType = iota + 1
	AWAITING_PAYMENT
	FAILED
	PAYED
	CANCELLED
)

type Order struct {
	ID     int64
	State  StateType
	Items  []*Item
	UserId int64
}
