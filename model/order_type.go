package model

const (
	HOLD    = OrderType("Hold")
	SUPPORT = OrderType("Support")
	MOVE    = OrderType("Move")
)

type OrderType string
