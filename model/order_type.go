package model

const (
	HOLD    = OrderType("hold")
	SUPPORT = OrderType("support")
	MOVE    = OrderType("move")
)

type OrderType string
