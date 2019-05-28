package model

const (
	CONVOY        = OrderType("Convoy")
	HOLD          = OrderType("Hold")
	MOVE          = OrderType("Move")
	MOVEVIACONVOY = OrderType("MoveConvoy")
	SUPPORT       = OrderType("Support")
	RETREAT       = OrderType("Retreat")
	BUILD         = OrderType("Build")
	DESTROY       = OrderType("Destroy")
)

type OrderType string
