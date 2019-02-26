package model

import (
	"errors"
	"time"
)

type Game struct {
	Id             int        `json: "id"`
	Title          string     `json: "title"`
	StartedAt      *time.Time `json: "started_at"`
	GameYear       *time.Time `json: "game_year"`
	Phase          int        `json: "phase"`
	PhaseEnd       string     `json: "phase_end"`
	OrdersInterval int        `json: "orders_interval"`
	GameBoard      GameBoard  `json: "game_squares,omitempty"`
}

// DrawGameBoard() Fills in an empty game.GameBoard from data loaded from territory and piece rows in the db
// Its reads from the db and can be used at any time to get the current game board state
func (g *Game) DrawGameBoard(territoryRows []TerritoryRow, pieceRows []PieceRow) {
	var gb = GameBoard{}
	var unit = &Unit{}

	for _, row := range territoryRows {
		units := make([]Unit, 0)
		for _, pr := range pieceRows {
			if row.Country == pr.Country {
				unit.UnitType = pr.UnitType
				unit.Owner = pr.Owner
				unit.PieceId = pr.Id
				units = append(units, *unit)
			}
		}
		gb[row.Country] = GameSquareData{Owner: row.Owner, Units: units, TerritoryId: row.Id}
	}
	g.GameBoard = gb
}

// NewGame() creates a new Game model
func (g *Game) NewGameBoard() {
	var gb = GameBoard{

		EDINBURGH: GameSquareData{Owner: ENGLAND, Units: []Unit{{UnitType: NAVY, Owner: ENGLAND}}},
		LONDON:    GameSquareData{Owner: ENGLAND, Units: []Unit{{UnitType: NAVY, Owner: ENGLAND}}},
		LIVERPOOL: GameSquareData{Owner: ENGLAND, Units: []Unit{{UnitType: ARMY, Owner: ENGLAND}}},
		WALES:     GameSquareData{Owner: ENGLAND},
		YORKSHIRE: GameSquareData{Owner: ENGLAND},
		CLYDE:     GameSquareData{Owner: ENGLAND},

		KIEL:    GameSquareData{Owner: GERMANY, Units: []Unit{{UnitType: NAVY, Owner: GERMANY}}},
		BERLIN:  GameSquareData{Owner: GERMANY, Units: []Unit{{UnitType: ARMY, Owner: GERMANY}}},
		MUNICH:  GameSquareData{Owner: GERMANY, Units: []Unit{{UnitType: ARMY, Owner: GERMANY}}},
		RUHR:    GameSquareData{Owner: GERMANY},
		SILESIA: GameSquareData{Owner: GERMANY},
		PRUSSIA: GameSquareData{Owner: GERMANY},

		BREST:      GameSquareData{Owner: FRANCE, Units: []Unit{{UnitType: NAVY, Owner: FRANCE}}},
		PARIS:      GameSquareData{Owner: FRANCE, Units: []Unit{{UnitType: ARMY, Owner: FRANCE}}},
		MARSEILLES: GameSquareData{Owner: FRANCE, Units: []Unit{{UnitType: ARMY, Owner: FRANCE}}},
		GASCONY:    GameSquareData{Owner: FRANCE},
		BURGUNDY:   GameSquareData{Owner: FRANCE},
		PICARDY:    GameSquareData{Owner: FRANCE},

		NAPLES:   GameSquareData{Owner: ITALY, Units: []Unit{{UnitType: NAVY, Owner: ITALY}}},
		ROME:     GameSquareData{Owner: ITALY, Units: []Unit{{UnitType: ARMY, Owner: ITALY}}},
		VENICE:   GameSquareData{Owner: ITALY, Units: []Unit{{UnitType: ARMY, Owner: ITALY}}},
		TUSCANY:  GameSquareData{Owner: ITALY},
		PIEDMONT: GameSquareData{Owner: ITALY},
		APULIA:   GameSquareData{Owner: ITALY},

		VIENNA:   GameSquareData{Owner: AUSTRIA_HUNGARY, Units: []Unit{{UnitType: ARMY, Owner: AUSTRIA_HUNGARY}}},
		TRIESTE:  GameSquareData{Owner: AUSTRIA_HUNGARY, Units: []Unit{{UnitType: NAVY, Owner: AUSTRIA_HUNGARY}}},
		BUDAPEST: GameSquareData{Owner: AUSTRIA_HUNGARY, Units: []Unit{{UnitType: ARMY, Owner: AUSTRIA_HUNGARY}}},
		GALICIA:  GameSquareData{Owner: AUSTRIA_HUNGARY},
		TYROLIA:  GameSquareData{Owner: AUSTRIA_HUNGARY},
		BOHEMIA:  GameSquareData{Owner: AUSTRIA_HUNGARY},

		CONSTANTINOPLE: GameSquareData{Owner: TURKEY, Units: []Unit{{UnitType: ARMY, Owner: TURKEY}}},
		ANKARA:         GameSquareData{Owner: TURKEY, Units: []Unit{{UnitType: NAVY, Owner: TURKEY}}},
		SMYRNA:         GameSquareData{Owner: TURKEY, Units: []Unit{{UnitType: ARMY, Owner: TURKEY}}},
		ARMENIA:        GameSquareData{Owner: TURKEY},
		SYRIA:          GameSquareData{Owner: TURKEY},

		ST_PETERSBURG_SOUTH_COAST: GameSquareData{Owner: RUSSIA, Units: []Unit{{UnitType: NAVY, Owner: RUSSIA}}},
		SEVASTOPOL:                GameSquareData{Owner: RUSSIA, Units: []Unit{{UnitType: NAVY, Owner: RUSSIA}}},
		MOSCOW:                    GameSquareData{Owner: RUSSIA, Units: []Unit{{UnitType: ARMY, Owner: RUSSIA}}},
		WARSAW:                    GameSquareData{Owner: RUSSIA, Units: []Unit{{UnitType: ARMY, Owner: RUSSIA}}},
		UKRAINE:                   GameSquareData{Owner: RUSSIA},
		ST_PETERSBURG:             GameSquareData{Owner: RUSSIA},
		ST_PETERSBURG_NORTH_COAST: GameSquareData{Owner: RUSSIA},
		LIVONIA:                   GameSquareData{Owner: RUSSIA},
		FINLAND:                   GameSquareData{Owner: RUSSIA},

		AEGEAN_SEA:            GameSquareData{},
		ADRIATIC_SEA:          GameSquareData{},
		SERBIA:                GameSquareData{},
		ALBANIA:               GameSquareData{},
		GREECE:                GameSquareData{},
		ROMANIA:               GameSquareData{},
		BULGARIA:              GameSquareData{},
		BULGARIA_EAST_COAST:   GameSquareData{},
		BULGARIA_SOUTH_COAST:  GameSquareData{},
		AEGEAN_SESA:           GameSquareData{},
		EASTERN_MEDITERRANEAN: GameSquareData{},
		NORTH_ATLANTIC_OCEAN:  GameSquareData{},
		IRISH_SEA:             GameSquareData{},
		ENGLISH_CHANNEL:       GameSquareData{},
		NORTH_SEA:             GameSquareData{},
		NORWEGIAN_SEA:         GameSquareData{},
		MID_ATLANTIC_OCEAN:    GameSquareData{},
		BELGIUM:               GameSquareData{},
		HOLLAND:               GameSquareData{},
		SPAIN:                 GameSquareData{},
		SPAIN_NORTH_COAST:     GameSquareData{},
		SPAIN_SOUTH_COAST:     GameSquareData{},
		PORTUGAL:              GameSquareData{},
		GULF_OF_LYON:          GameSquareData{},
		WESTERN_MEDITERRANEAN: GameSquareData{},
		NORTH_AFRICA:          GameSquareData{},
		TUNIS:                 GameSquareData{},
		TYRRHENIAN_SEA:        GameSquareData{},
		IONIAN_SEA:            GameSquareData{},
		BLACK_SEA:             GameSquareData{},
		SWEDEN:                GameSquareData{},
		NORWAY:                GameSquareData{},
		GULF_OF_BOTHNIA:       GameSquareData{},
		BARRENTS_SEA:          GameSquareData{},
		BALTIC_SEA:            GameSquareData{},
		DENMARK:               GameSquareData{},
		SKAGERRAK:             GameSquareData{},
		HELGOLAND_BIGHT:       GameSquareData{},
	}

	g.GameBoard = gb
}

func (g *Game) ValidatePhaseUpdate(updateToPhase int) error {
	if g.Phase >= updateToPhase {
		return errors.New("The Phase has already been updated")
	}
	return nil
}
