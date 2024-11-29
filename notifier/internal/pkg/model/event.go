package model

type StateType string
type SKUType uint32

const (
	NEW              StateType = "NEW"
	AWAITING_PAYMENT           = "AWAITING_PAYMENT"
	FAILED                     = "FAILED"
	PAYED                      = "PAYED"
	CANCELLED                  = "CANCELLED"
)

type OrderEvent struct {
	ID     int64
	State  StateType
	Items  []*Item
	UserId int64
}

type Item struct {
	SKU   SKUType
	Count uint32
}
