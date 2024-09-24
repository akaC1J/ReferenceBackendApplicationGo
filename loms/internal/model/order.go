package model

type StateType string

const (
	NEW              StateType = "NEW"
	AWAITING_PAYMENT           = "AWAITING_PAYMENT"
	FAILED                     = "FAILED"
	PAYED                      = "PAYED"
	CANCELLED                  = "CANCELLED"
)

type Order struct {
	ID     int64
	State  StateType
	Items  []*Item
	UserId int64
}
