package model

const (
	CONVOY        = OrderType("Convoy")
	HOLD          = OrderType("Hold")
	MOVE          = OrderType("Move")
	MOVEVIACONVOY = OrderType("MoveConvoy")
	SUPPORT       = OrderType("Support")
)

type OrderType string
