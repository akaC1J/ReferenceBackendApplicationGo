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

type Order struct {
	ID     int64
	state  StateType
	Items  []*Item
	UserId int64
}

func (order *Order) SetState(state StateType) error {
	availableStates := map[StateType]struct{}{
		NEW:              {},
		AWAITING_PAYMENT: {},
		FAILED:           {},
		PAYED:            {},
		CANCELLED:        {},
	}
	state = StateType(strings.ToUpper(string(state)))
	if _, ok := availableStates[state]; !ok {
		return fmt.Errorf("invalid state: %s", state)
	}
	order.state = state
	return nil
}

func (order *Order) State() StateType {
	return order.state
}
