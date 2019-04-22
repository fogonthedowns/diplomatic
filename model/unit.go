package model

type Unit struct {
	PieceId     int64    `json:"piece_id,omitempty"`
	UnitType    UnitType `json:"unit_type"`
	Owner       Country  `json:"owner"`
	WillRetreat bool     `json:"will_retreat,omitempty"`
}

const (
	ARMY = UnitType("Army")
	NAVY = UnitType("Navy")
)

type UnitType string
