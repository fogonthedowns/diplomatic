package model

const (
	CONVOY        = OrderType("Convoy")
	HOLD          = OrderType("Hold")
	MOVE          = OrderType("Move")
	MOVEVIACONVOY = OrderType("MoveViaConvoy")
	SUPPORT       = OrderType("Support")
)

type OrderType string
