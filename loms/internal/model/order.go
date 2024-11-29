package model

import (
	"fmt"
	"strings"
)

type StateType string

const (
	NEW              StateType = "NEW"
	AWAITING_PAYMENT           = "AWAITING_PAYMENT"
	FAILED                     = "FAILED"
	PAYED                      = "PAYED"
	CANCELLED                  = "CANCELLED"
)

var availableStates = map[StateType]struct{}{
	NEW:              {},
	AWAITING_PAYMENT: {},
	FAILED:           {},
	PAYED:            {},
	CANCELLED:        {},
}

type Order struct {
	ID     int64
	State  StateType
	Items  []*Item
	UserId int64
}

func (order *Order) SetState(state StateType) error {
	state = StateType(strings.ToUpper(string(state)))
	if _, ok := availableStates[state]; !ok {
		return fmt.Errorf("invalid State: %s", state)
	}
	order.State = state
	return nil
}
