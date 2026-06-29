package oracle

import (
	"errors"
)

type MockOracle struct {
	prices     map[int64]int64
	failOnce   map[int64]bool // simulate transient failure
	alwaysFail bool           // simulate oracle being completely down
	returnZero map[int64]bool // simulate bad data (zero price)
	returnNeg  map[int64]bool // simulate bad data (negative price)
	CallCount  map[int64]int  // track how many times each item was fetched
}

func NewMock() *MockOracle {
	return &MockOracle{
		prices:     make(map[int64]int64),
		failOnce:   make(map[int64]bool),
		returnZero: make(map[int64]bool),
		returnNeg:  make(map[int64]bool),
		CallCount:  make(map[int64]int),
	}
}

func (m *MockOracle) SetPrice(itemID int64, price int64) {
	m.prices[itemID] = price
}

func (m *MockOracle) SetFailOnce(itemID int64) {
	m.failOnce[itemID] = true
}

func (m *MockOracle) SetAlwaysFail() {
	m.alwaysFail = true
}

func (m *MockOracle) SetReturnZero(itemID int64) {
	m.returnZero[itemID] = true
}

func (m *MockOracle) SetReturnNegative(itemID int64) {
	m.returnNeg[itemID] = true
}

func (m *MockOracle) GetBasePrice(itemID int64) (int64, error) {
	m.CallCount[itemID]++

	if m.alwaysFail {
		return 0, errors.New("oracle is down")
	}

	if m.failOnce[itemID] {
		delete(m.failOnce, itemID)
		return 0, errors.New("transient oracle failure")
	}

	if m.returnZero[itemID] {
		return 0, nil
	}

	if m.returnNeg[itemID] {
		return -100, nil
	}

	price, ok := m.prices[itemID]
	if !ok {
		return 0, errors.New("price not found")
	}

	return price, nil
}
