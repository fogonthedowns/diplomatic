package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

const (
	AEGEAN_SEA                = Territory("AEG")
	CLYDE                     = Territory("CLY")
	NORTH_ATLANTIC_OCEAN      = Territory("NAO")
	IRISH_SEA                 = Territory("IHS")
	ENGLISH_CHANNEL           = Territory("ENC")
	EDINBURGH                 = Territory("EDB")
	YORKSHIRE                 = Territory("YRK")
	LONDON                    = Territory("LON")
	WALES                     = Territory("WLS")
	LIVERPOOL                 = Territory("LIV")
	NORTH_SEA                 = Territory("NHS")
	NORWEGIAN_SEA             = Territory("NOS")
	MID_ATLANTIC_OCEAN        = Territory("MAO")
	BREST                     = Territory("BST")
	PICARDY                   = Territory("PIC")
	PARIS                     = Territory("PRS")
	BELGIUM                   = Territory("BGM")
	HOLLAND                   = Territory("HOL")
	GASCONY                   = Territory("GAS")
	BURGUNDY                  = Territory("BRG")
	MARSEILLES                = Territory("MAR")
	SPAIN                     = Territory("SPA")
	SPAIN_NORTH_COAST         = Territory("SNC")
	SPAIN_SOUTH_COAST         = Territory("SSC")
	PORTUGAL                  = Territory("PRT")
	GULF_OF_LYON              = Territory("GOL")
	WESTERN_MEDITERRANEAN     = Territory("WMD")
	NORTH_AFRICA              = Territory("NAK")
	TUNIS                     = Territory("TUN")
	TYRRHENIAN_SEA            = Territory("TYR")
	IONIAN_SEA                = Territory("ION")
	PIEDMONT                  = Territory("PDT")
	VENICE                    = Territory("VEN")
	TUSCANY                   = Territory("TUS")
	ROME                      = Territory("ROM")
	APULIA                    = Territory("APU")
	NAPLES                    = Territory("NAP")
	ADRIATIC_SEA              = Territory("ADR")
	TYROLIA                   = Territory("TYA")
	BOHEMIA                   = Territory("BOH")
	VIENNA                    = Territory("VNA")
	TRIESTE                   = Territory("TTE")
	BUDAPEST                  = Territory("BUD")
	SERBIA                    = Territory("SBA")
	ALBANIA                   = Territory("ALB")
	GREECE                    = Territory("GRC")
	ROMANIA                   = Territory("RMA")
	BULGARIA                  = Territory("BUL")
	BULGARIA_EAST_COAST       = Territory("BUE")
	BULGARIA_SOUTH_COAST      = Territory("BUS")
	AEGEAN_SESA               = Territory("AGS")
	EASTERN_MEDITERRANEAN     = Territory("EMD")
	CONSTANTINOPLE            = Territory("CON")
	ANKARA                    = Territory("ANK")
	SMYRNA                    = Territory("SMY")
	ARMENIA                   = Territory("ARM")
	SYRIA                     = Territory("SYR")
	BLACK_SEA                 = Territory("BLA")
	SEVASTOPOL                = Territory("SEV")
	MOSCOW                    = Territory("MOS")
	UKRAINE                   = Territory("UKR")
	ST_PETERSBURG             = Territory("STP")
	ST_PETERSBURG_NORTH_COAST = Territory("SPN")
	ST_PETERSBURG_SOUTH_COAST = Territory("SPS")
	LIVONIA                   = Territory("LVA")
	WARSAW                    = Territory("WAR")
	FINLAND                   = Territory("FIN")
	SWEDEN                    = Territory("SWE")
	NORWAY                    = Territory("NOR")
	GULF_OF_BOTHNIA           = Territory("GOB")
	BARRENTS_SEA              = Territory("BAS")
	BALTIC_SEA                = Territory("BAL")
	PRUSSIA                   = Territory("PRU")
	SILESIA                   = Territory("SIL")
	GALICIA                   = Territory("GAL")
	DENMARK                   = Territory("DEN")
	SKAGERRAK                 = Territory("SKA")
	HELGOLAND_BIGHT           = Territory("HEL")
	KIEL                      = Territory("KIE")
	BERLIN                    = Territory("BER")
	MUNICH                    = Territory("MUN")
	RUHR                      = Territory("RUH")
	BLANK                     = Territory("NONE")
)

var russiaEdgeCases = []Territory{ST_PETERSBURG, ST_PETERSBURG_NORTH_COAST, ST_PETERSBURG_SOUTH_COAST}
var bulgariaEdgeCases = []Territory{BULGARIA_SOUTH_COAST, BULGARIA_EAST_COAST, BULGARIA}
var spainEdgeCases = []Territory{SPAIN_NORTH_COAST, SPAIN_SOUTH_COAST, SPAIN}

// TODO
// mapOfCenters
// mapOfLandOrSea

type Territory string
type TerritoryMoves map[Territory][]*Move

func (t Territory) RussiaEdgeCase() bool {
	for _, c := range russiaEdgeCases {
		if t == c {
			return true

		}
	}
	return false
}

func (t Territory) SpainEdgeCase() bool {

	for _, c := range spainEdgeCases {
		if c == t {
			return true
		}
	}
	return false
}

func (t Territory) BulgariaEdgeCase() bool {
	for _, c := range bulgariaEdgeCases {
		if c == t {
			return true
		}
	}
	return false
}

func (tm TerritoryMoves) Uncontested(key Territory) bool {
	for _, c := range russiaEdgeCases {
		if key == c {
			return len(tm[ST_PETERSBURG])+len(tm[ST_PETERSBURG_NORTH_COAST])+len(tm[ST_PETERSBURG_SOUTH_COAST]) <= 1
		}
	}

	for _, c := range spainEdgeCases {
		if key == c {
			return len(tm[SPAIN_NORTH_COAST])+len(tm[SPAIN_SOUTH_COAST])+len(tm[SPAIN]) <= 1
		}
	}

	for _, c := range bulgariaEdgeCases {
		if key == c {
			return len(tm[BULGARIA_SOUTH_COAST])+len(tm[BULGARIA_EAST_COAST])+len(tm[BULGARIA]) <= 1
		}
	}
	return len(tm[key]) <= 1
}

// ResolveConflicts() resolves conflicts of any territory that has two Move objects assoicated with it.
// These objects represent a move conflict.
// Resolves a conflict using move.MovePower

func (tm TerritoryMoves) ResolveConflicts(moves *Moves) {
	for key, value := range tm {
		var lastSeen = 0
		var lastSeenMove *Move
		if len(tm[key]) >= 2 {
			value = Sort(value)
			for index, mm := range value {
				if mm.MovePower > 0 && mm.MovePower > lastSeen {
					if lastSeenMove != nil {
						lastSeenMove.DislodgeIfHold(moves, lastSeenMove)
					}
					value[index].MovePieceForward()

				}

				if mm.MovePower == lastSeen {
					if lastSeenMove != nil {
						lastSeenMove.BouncePiece()
					}
					value[index].BouncePiece()
				}

				if mm.MovePower < lastSeen {
					value[index].BouncePiece()
					value[index].DislodgeIfHold(moves, lastSeenMove)
				}

				lastSeen = mm.MovePower
				lastSeenMove = value[index]
			}
		}
	}

}

// Sort() sorts the Move Values by MovePower
// The values are sorted to prevent a bug where
// Disloged is never set
// because of a mismatched comparison
// such as 2, 0, 1 where - 1 is not disloged bc it is greater than 0.
func Sort(v []*Move) (nv []*Move) {
	var last *int
	for _, value := range v {
		if last == nil {
			nv = append(nv, value)
		} else {
			if value.MovePower > *last {
				nv = append([]*Move{value}, nv...)
			} else {
				nv = append(nv, value)
			}
		}
		last = &value.MovePower
	}
	return nv
}

// validSeaMoves defines a map of Valid moves for Navy Units
var validSeaMoves = map[Territory][]Territory{
	NORWEGIAN_SEA:             []Territory{NORTH_ATLANTIC_OCEAN, BARRENTS_SEA, NORWAY, NORTH_SEA},
	BARRENTS_SEA:              []Territory{ST_PETERSBURG_NORTH_COAST, NORWAY, NORWEGIAN_SEA},
	NORTH_ATLANTIC_OCEAN:      []Territory{NORWEGIAN_SEA, CLYDE, IRISH_SEA, MID_ATLANTIC_OCEAN, LIVERPOOL},
	IRISH_SEA:                 []Territory{WALES, LIVERPOOL, NORTH_ATLANTIC_OCEAN, MID_ATLANTIC_OCEAN, ENGLISH_CHANNEL},
	ENGLISH_CHANNEL:           []Territory{IRISH_SEA, MID_ATLANTIC_OCEAN, BREST, PICARDY, BELGIUM, LONDON, WALES, NORTH_SEA},
	NORTH_SEA:                 []Territory{EDINBURGH, YORKSHIRE, LONDON, ENGLISH_CHANNEL, BELGIUM, HOLLAND, HELGOLAND_BIGHT, DENMARK, SKAGERRAK, NORWAY, NORWEGIAN_SEA},
	HELGOLAND_BIGHT:           []Territory{KIEL, HOLLAND, DENMARK, NORTH_SEA},
	SKAGERRAK:                 []Territory{NORTH_SEA, NORWAY, SWEDEN, DENMARK},                               // TODO VERIFY BALTIC_SEA
	BALTIC_SEA:                []Territory{DENMARK, SWEDEN, GULF_OF_BOTHNIA, LIVONIA, PRUSSIA, BERLIN, KIEL}, // TODO VERIFY SKAGERRAK
	GULF_OF_BOTHNIA:           []Territory{SWEDEN, FINLAND, ST_PETERSBURG_SOUTH_COAST, LIVONIA, BALTIC_SEA},
	MID_ATLANTIC_OCEAN:        []Territory{NORTH_ATLANTIC_OCEAN, IRISH_SEA, ENGLISH_CHANNEL, BREST, GASCONY, SPAIN_NORTH_COAST, PORTUGAL, NORTH_AFRICA, WESTERN_MEDITERRANEAN},
	WESTERN_MEDITERRANEAN:     []Territory{SPAIN_SOUTH_COAST, GULF_OF_LYON, TYRRHENIAN_SEA, TUNIS, NORTH_AFRICA, MID_ATLANTIC_OCEAN},
	GULF_OF_LYON:              []Territory{MARSEILLES, TYRRHENIAN_SEA, WESTERN_MEDITERRANEAN, SPAIN_SOUTH_COAST, PIEDMONT},
	TYRRHENIAN_SEA:            []Territory{GULF_OF_LYON, TUSCANY, ROME, NAPLES, IONIAN_SEA, TUNIS, WESTERN_MEDITERRANEAN},
	IONIAN_SEA:                []Territory{APULIA, TYRRHENIAN_SEA, NAPLES, ADRIATIC_SEA, ALBANIA, GREECE, AEGEAN_SEA, EASTERN_MEDITERRANEAN},
	ADRIATIC_SEA:              []Territory{VENICE, TRIESTE, ALBANIA, IONIAN_SEA, APULIA},
	AEGEAN_SEA:                []Territory{GREECE, BULGARIA_SOUTH_COAST, CONSTANTINOPLE, SMYRNA, EASTERN_MEDITERRANEAN, IONIAN_SEA},
	BLACK_SEA:                 []Territory{SEVASTOPOL, ARMENIA, ANKARA, CONSTANTINOPLE, BULGARIA_EAST_COAST, ROMANIA},
	SEVASTOPOL:                []Territory{BLACK_SEA, ROMANIA, ARMENIA},
	ARMENIA:                   []Territory{SEVASTOPOL, ANKARA, BLACK_SEA},
	ANKARA:                    []Territory{ARMENIA, BLACK_SEA, CONSTANTINOPLE},
	CONSTANTINOPLE:            []Territory{ANKARA, BLACK_SEA, BULGARIA_EAST_COAST, BULGARIA_SOUTH_COAST, AEGEAN_SEA},
	BULGARIA_EAST_COAST:       []Territory{CONSTANTINOPLE, BLACK_SEA, ROMANIA},
	ROMANIA:                   []Territory{SEVASTOPOL, BLACK_SEA, BULGARIA_EAST_COAST},
	BULGARIA_SOUTH_COAST:      []Territory{CONSTANTINOPLE, AEGEAN_SEA, GREECE},
	SMYRNA:                    []Territory{CONSTANTINOPLE, SYRIA, EASTERN_MEDITERRANEAN, AEGEAN_SEA},
	SYRIA:                     []Territory{SMYRNA, EASTERN_MEDITERRANEAN},
	GREECE:                    []Territory{AEGEAN_SEA, IONIAN_SEA, ALBANIA, BULGARIA_SOUTH_COAST},
	ALBANIA:                   []Territory{GREECE, TRIESTE, IONIAN_SEA, ADRIATIC_SEA},
	TRIESTE:                   []Territory{VENICE, ADRIATIC_SEA, ALBANIA},
	VENICE:                    []Territory{TRIESTE, ADRIATIC_SEA, APULIA},
	APULIA:                    []Territory{ADRIATIC_SEA, VENICE, NAPLES, IONIAN_SEA},
	NAPLES:                    []Territory{APULIA, IONIAN_SEA, TYRRHENIAN_SEA, ROME},
	ROME:                      []Territory{NAPLES, TUSCANY, TYRRHENIAN_SEA},
	TUSCANY:                   []Territory{ROME, PIEDMONT, GULF_OF_LYON, TYRRHENIAN_SEA},
	PIEDMONT:                  []Territory{MARSEILLES, TUSCANY, GULF_OF_LYON},
	MARSEILLES:                []Territory{PIEDMONT, SPAIN_SOUTH_COAST, GULF_OF_LYON},
	SPAIN_SOUTH_COAST:         []Territory{MARSEILLES, GULF_OF_LYON, WESTERN_MEDITERRANEAN, MID_ATLANTIC_OCEAN, PORTUGAL}, // TODO verify MID_ATLANTIC_OCEAN/PORTUGAL
	NORTH_AFRICA:              []Territory{WESTERN_MEDITERRANEAN, TUNIS, MID_ATLANTIC_OCEAN},
	TUNIS:                     []Territory{IONIAN_SEA, TYRRHENIAN_SEA, WESTERN_MEDITERRANEAN, NORTH_AFRICA},
	PORTUGAL:                  []Territory{MID_ATLANTIC_OCEAN, SPAIN_SOUTH_COAST, SPAIN_NORTH_COAST}, // TODO verify SPAIN_SOUTH_COAST
	SPAIN_NORTH_COAST:         []Territory{MID_ATLANTIC_OCEAN, GASCONY, PORTUGAL},                    // TODO Verify north/south coasts do not border
	GASCONY:                   []Territory{SPAIN_NORTH_COAST, BREST, MID_ATLANTIC_OCEAN},
	BREST:                     []Territory{GASCONY, MID_ATLANTIC_OCEAN, ENGLISH_CHANNEL, PICARDY},
	PICARDY:                   []Territory{BREST, ENGLISH_CHANNEL, BELGIUM},
	BELGIUM:                   []Territory{PICARDY, ENGLISH_CHANNEL, NORTH_SEA, HOLLAND},
	HOLLAND:                   []Territory{BELGIUM, NORTH_SEA, HELGOLAND_BIGHT, KIEL},
	KIEL:                      []Territory{HELGOLAND_BIGHT, DENMARK, BALTIC_SEA, BERLIN, HOLLAND},
	BERLIN:                    []Territory{KIEL, BALTIC_SEA, PRUSSIA},
	PRUSSIA:                   []Territory{BERLIN, BALTIC_SEA, LIVONIA},
	LIVONIA:                   []Territory{PRUSSIA, BALTIC_SEA, GULF_OF_BOTHNIA, ST_PETERSBURG_SOUTH_COAST},
	ST_PETERSBURG_SOUTH_COAST: []Territory{LIVONIA, GULF_OF_BOTHNIA, FINLAND},
	FINLAND:                   []Territory{ST_PETERSBURG_SOUTH_COAST, SWEDEN},
	SWEDEN:                    []Territory{FINLAND, GULF_OF_BOTHNIA, BALTIC_SEA, SKAGERRAK, NORWAY},
	NORWAY:                    []Territory{SWEDEN, SKAGERRAK, NORTH_SEA, NORWEGIAN_SEA, BARRENTS_SEA, ST_PETERSBURG_NORTH_COAST},
	ST_PETERSBURG_NORTH_COAST: []Territory{NORWAY, BARRENTS_SEA},
	EDINBURGH:                 []Territory{CLYDE, YORKSHIRE, NORTH_SEA, NORWEGIAN_SEA},
	CLYDE:                     []Territory{EDINBURGH, NORWEGIAN_SEA, NORTH_ATLANTIC_OCEAN, LIVERPOOL},
	LIVERPOOL:                 []Territory{CLYDE, WALES, IRISH_SEA, NORTH_ATLANTIC_OCEAN},
	WALES:                     []Territory{LIVERPOOL, ENGLISH_CHANNEL, IRISH_SEA, LONDON},
	LONDON:                    []Territory{WALES, ENGLISH_CHANNEL, NORTH_SEA, YORKSHIRE},
	YORKSHIRE:                 []Territory{LONDON, EDINBURGH, NORTH_SEA},
}

// validLandMovess defines a map of Valid moves for Navy Units
// The keys of this are also available land territories
var validLandMoves = map[Territory][]Territory{
	CLYDE:          []Territory{EDINBURGH, LIVERPOOL},
	EDINBURGH:      []Territory{CLYDE, YORKSHIRE},
	LIVERPOOL:      []Territory{CLYDE, EDINBURGH, YORKSHIRE, WALES},
	WALES:          []Territory{LIVERPOOL, LONDON, YORKSHIRE},
	YORKSHIRE:      []Territory{LIVERPOOL, LONDON, WALES, EDINBURGH},
	BREST:          []Territory{PICARDY, PARIS, GASCONY},
	PICARDY:        []Territory{BREST, BELGIUM, BURGUNDY, PARIS},
	PARIS:          []Territory{GASCONY, BREST, PICARDY, BURGUNDY},
	BURGUNDY:       []Territory{MARSEILLES, GASCONY, PARIS, BELGIUM, RUHR, MUNICH},
	GASCONY:        []Territory{BREST, PARIS, BURGUNDY, MARSEILLES, SPAIN},
	MARSEILLES:     []Territory{SPAIN, GASCONY, BURGUNDY, PIEDMONT},
	SPAIN:          []Territory{MARSEILLES, GASCONY, PORTUGAL},
	PORTUGAL:       []Territory{SPAIN},
	BELGIUM:        []Territory{PICARDY, BURGUNDY, RUHR, HOLLAND},
	HOLLAND:        []Territory{BELGIUM, RUHR, KIEL},
	DENMARK:        []Territory{KIEL},
	KIEL:           []Territory{HOLLAND, DENMARK, BERLIN, MUNICH, RUHR},
	RUHR:           []Territory{BELGIUM, HOLLAND, KIEL, MUNICH, BURGUNDY},
	MUNICH:         []Territory{BURGUNDY, RUHR, TYROLIA, BOHEMIA, SILESIA, BERLIN, KIEL},
	BOHEMIA:        []Territory{MUNICH, SILESIA, GALICIA, VIENNA, TYROLIA},
	SILESIA:        []Territory{MUNICH, BOHEMIA, BERLIN, PRUSSIA, WARSAW, GALICIA},
	PRUSSIA:        []Territory{BERLIN, LIVONIA, WARSAW, SILESIA},
	LIVONIA:        []Territory{PRUSSIA, WARSAW, MOSCOW, ST_PETERSBURG},
	ST_PETERSBURG:  []Territory{LIVONIA, MOSCOW, FINLAND},
	MOSCOW:         []Territory{ST_PETERSBURG, LIVONIA, WARSAW, UKRAINE, SEVASTOPOL},
	WARSAW:         []Territory{LIVONIA, MOSCOW, UKRAINE, GALICIA, SILESIA, PRUSSIA},
	UKRAINE:        []Territory{WARSAW, MOSCOW, SEVASTOPOL, ROMANIA, GALICIA},
	SEVASTOPOL:     []Territory{UKRAINE, ARMENIA, MOSCOW, ROMANIA},
	GALICIA:        []Territory{BOHEMIA, SILESIA, WARSAW, UKRAINE, ROMANIA, BUDAPEST, VIENNA},
	ROMANIA:        []Territory{BUDAPEST, GALICIA, UKRAINE, SEVASTOPOL, BULGARIA, SERBIA},
	BULGARIA:       []Territory{GREECE, SERBIA, ROMANIA, CONSTANTINOPLE},
	BUDAPEST:       []Territory{GALICIA, ROMANIA, SERBIA, TRIESTE, VIENNA},
	SERBIA:         []Territory{ROMANIA, BULGARIA, GREECE, ALBANIA, TRIESTE, BUDAPEST},
	ALBANIA:        []Territory{TRIESTE, SERBIA, GREECE},
	GREECE:         []Territory{ALBANIA, SERBIA, BUDAPEST},
	VIENNA:         []Territory{BOHEMIA, GALICIA, BUDAPEST, TRIESTE, TYROLIA},
	TRIESTE:        []Territory{SERBIA, ALBANIA, VENICE, TYROLIA, VIENNA, BUDAPEST},
	TYROLIA:        []Territory{MUNICH, BOHEMIA, VIENNA, TRIESTE, VENICE, PIEDMONT},
	VENICE:         []Territory{TYROLIA, TRIESTE, PIEDMONT, TUSCANY, APULIA, ROME},
	PIEDMONT:       []Territory{MARSEILLES, TYROLIA, VENICE, TUSCANY},
	TUSCANY:        []Territory{PIEDMONT, VENICE, ROME},
	ROME:           []Territory{TUSCANY, VENICE, APULIA, NAPLES},
	APULIA:         []Territory{VENICE, ROME, NAPLES},
	NAPLES:         []Territory{APULIA, ROME},
	TUNIS:          []Territory{NORTH_AFRICA},
	NORTH_AFRICA:   []Territory{TUNIS},
	SYRIA:          []Territory{ARMENIA, SMYRNA},
	ARMENIA:        []Territory{SYRIA, SMYRNA, ANKARA, SEVASTOPOL},
	CONSTANTINOPLE: []Territory{ANKARA, SMYRNA, BULGARIA},
	SMYRNA:         []Territory{CONSTANTINOPLE, ANKARA, ARMENIA, SYRIA},
	ANKARA:         []Territory{CONSTANTINOPLE, SMYRNA, ARMENIA},
	FINLAND:        []Territory{ST_PETERSBURG, NORWAY, SWEDEN},
	SWEDEN:         []Territory{NORWAY, FINLAND},
	NORWAY:         []Territory{SWEDEN, FINLAND, ST_PETERSBURG},
}

var validLandSupportMoves = map[Territory][]Territory{
	CLYDE:          []Territory{EDINBURGH, LIVERPOOL},
	EDINBURGH:      []Territory{CLYDE, YORKSHIRE},
	LIVERPOOL:      []Territory{CLYDE, EDINBURGH, YORKSHIRE, WALES},
	WALES:          []Territory{LIVERPOOL, LONDON, YORKSHIRE},
	YORKSHIRE:      []Territory{LIVERPOOL, LONDON, WALES, EDINBURGH},
	BREST:          []Territory{PICARDY, PARIS, GASCONY},
	PICARDY:        []Territory{BREST, BELGIUM, BURGUNDY, PARIS},
	PARIS:          []Territory{GASCONY, BREST, PICARDY, BURGUNDY},
	BURGUNDY:       []Territory{MARSEILLES, GASCONY, PARIS, BELGIUM, RUHR, MUNICH},
	GASCONY:        []Territory{BREST, PARIS, BURGUNDY, MARSEILLES, SPAIN},
	MARSEILLES:     []Territory{SPAIN, GASCONY, BURGUNDY, PIEDMONT, SPAIN_SOUTH_COAST},
	SPAIN:          []Territory{MARSEILLES, GASCONY, PORTUGAL},
	PORTUGAL:       []Territory{SPAIN, SPAIN_NORTH_COAST, SPAIN_SOUTH_COAST},
	BELGIUM:        []Territory{PICARDY, BURGUNDY, RUHR, HOLLAND},
	HOLLAND:        []Territory{BELGIUM, RUHR, KIEL},
	DENMARK:        []Territory{KIEL},
	KIEL:           []Territory{HOLLAND, DENMARK, BERLIN, MUNICH, RUHR},
	RUHR:           []Territory{BELGIUM, HOLLAND, KIEL, MUNICH, BURGUNDY},
	MUNICH:         []Territory{BURGUNDY, RUHR, TYROLIA, BOHEMIA, SILESIA, BERLIN, KIEL},
	BOHEMIA:        []Territory{MUNICH, SILESIA, GALICIA, VIENNA, TYROLIA},
	SILESIA:        []Territory{MUNICH, BOHEMIA, BERLIN, PRUSSIA, WARSAW, GALICIA},
	PRUSSIA:        []Territory{BERLIN, LIVONIA, WARSAW, SILESIA},
	LIVONIA:        []Territory{PRUSSIA, WARSAW, MOSCOW, ST_PETERSBURG, ST_PETERSBURG_SOUTH_COAST},
	ST_PETERSBURG:  []Territory{LIVONIA, MOSCOW, FINLAND},
	MOSCOW:         []Territory{ST_PETERSBURG, LIVONIA, WARSAW, UKRAINE, SEVASTOPOL},
	WARSAW:         []Territory{LIVONIA, MOSCOW, UKRAINE, GALICIA, SILESIA, PRUSSIA},
	UKRAINE:        []Territory{WARSAW, MOSCOW, SEVASTOPOL, ROMANIA, GALICIA},
	SEVASTOPOL:     []Territory{UKRAINE, ARMENIA, MOSCOW, ROMANIA},
	GALICIA:        []Territory{BOHEMIA, SILESIA, WARSAW, UKRAINE, ROMANIA, BUDAPEST, VIENNA},
	ROMANIA:        []Territory{BUDAPEST, GALICIA, UKRAINE, SEVASTOPOL, BULGARIA, SERBIA, BULGARIA_EAST_COAST},
	BULGARIA:       []Territory{GREECE, SERBIA, ROMANIA, CONSTANTINOPLE},
	BUDAPEST:       []Territory{GALICIA, ROMANIA, SERBIA, TRIESTE, VIENNA},
	SERBIA:         []Territory{ROMANIA, BULGARIA, GREECE, ALBANIA, TRIESTE, BUDAPEST},
	ALBANIA:        []Territory{TRIESTE, SERBIA, GREECE},
	GREECE:         []Territory{ALBANIA, SERBIA, BUDAPEST, BULGARIA_SOUTH_COAST},
	VIENNA:         []Territory{BOHEMIA, GALICIA, BUDAPEST, TRIESTE, TYROLIA},
	TRIESTE:        []Territory{SERBIA, ALBANIA, VENICE, TYROLIA, VIENNA, BUDAPEST},
	TYROLIA:        []Territory{MUNICH, BOHEMIA, VIENNA, TRIESTE, VENICE, PIEDMONT},
	VENICE:         []Territory{TYROLIA, TRIESTE, PIEDMONT, TUSCANY, APULIA, ROME},
	PIEDMONT:       []Territory{MARSEILLES, TYROLIA, VENICE, TUSCANY},
	TUSCANY:        []Territory{PIEDMONT, VENICE, ROME},
	ROME:           []Territory{TUSCANY, VENICE, APULIA, NAPLES},
	APULIA:         []Territory{VENICE, ROME, NAPLES},
	NAPLES:         []Territory{APULIA, ROME},
	TUNIS:          []Territory{NORTH_AFRICA},
	NORTH_AFRICA:   []Territory{TUNIS},
	SYRIA:          []Territory{ARMENIA, SMYRNA},
	ARMENIA:        []Territory{SYRIA, SMYRNA, ANKARA, SEVASTOPOL},
	CONSTANTINOPLE: []Territory{ANKARA, SMYRNA, BULGARIA, BULGARIA_SOUTH_COAST, BULGARIA_EAST_COAST},
	SMYRNA:         []Territory{CONSTANTINOPLE, ANKARA, ARMENIA, SYRIA},
	ANKARA:         []Territory{CONSTANTINOPLE, SMYRNA, ARMENIA},
	FINLAND:        []Territory{ST_PETERSBURG, NORWAY, SWEDEN},
	SWEDEN:         []Territory{NORWAY, FINLAND},
	NORWAY:         []Territory{SWEDEN, FINLAND, ST_PETERSBURG, ST_PETERSBURG_NORTH_COAST},
}

// exclusiveSeaTerritories defines a list of Sea Exclusive territories
var exclusiveSeaTerritories = []Territory{
	NORWEGIAN_SEA,
	BARRENTS_SEA,
	NORTH_ATLANTIC_OCEAN,
	IRISH_SEA,
	ENGLISH_CHANNEL,
	NORTH_SEA,
	HELGOLAND_BIGHT,
	SKAGERRAK,
	BALTIC_SEA,
	GULF_OF_BOTHNIA,
	MID_ATLANTIC_OCEAN,
	WESTERN_MEDITERRANEAN,
	GULF_OF_LYON,
	TYRRHENIAN_SEA,
	IONIAN_SEA,
	ADRIATIC_SEA,
	AEGEAN_SEA,
	BLACK_SEA,
	BULGARIA_EAST_COAST,
	BULGARIA_SOUTH_COAST,
	SPAIN_SOUTH_COAST,
	SPAIN_NORTH_COAST,
	ST_PETERSBURG_SOUTH_COAST,
	ST_PETERSBURG_NORTH_COAST,
}

// ValidSeaMovement will return true if the checked territory
// is included inside of the mapOfBorders map
// uses the origional terriotry as they key
func (t *Territory) ValidSeaMovement(check Territory) bool {
	for _, borderTerritory := range validSeaMoves[*t] {
		if borderTerritory == check {
			return true
		}
	}
	return false
}

// ValidLandMovement will return true if the checked territory
// is included inside of the mapOfBorders map
// uses the origional terriotry as they key
func (t *Territory) ValidLandMovement(check Territory, mt OrderType) bool {
	if mt == SUPPORT {
		return validLandSupportMovement(check, *t)
	}
	for _, borderTerritory := range validLandMoves[*t] {
		if borderTerritory == check {
			return true
		}
	}
	return false
}

// validLandSupportMovement will return true if the checked territory
// is included inside of the mapOfBorders map
// uses the origional terriotry as they key
func validLandSupportMovement(check Territory, t Territory) bool {
	for _, borderTerritory := range validLandSupportMoves[t] {
		if borderTerritory == check {
			return true
		}
	}
	return false
}

// ValidConvoyBeginAndEnd will return true if the initial territory
// and the destination territory are both included as
// keys in validSeaMoves
func (t *Territory) ValidConvoyBeginAndEnd(check Territory) bool {
	if _, ok := validSeaMoves[*t]; ok {
		if _, endOk := validSeaMoves[check]; endOk {
			return true
		}
	}
	return false
}

var reTerritory = regexp.MustCompile(`^(|AEG|CLY|NAO|IHS|ENC|EDB|YRK|LON|WLS|LIV|NHS|NOS|MAO|BST|PIC|PRS|BGM|HOL|GAS|BRG|MAR|SPA|SNC|SSC||PRT|GOL|WMD|NAK|TUN|TYR|ION|PDT|VEN|TUS|ROM|APU|NAP|ADR|TYA|BOH|VNA|TTE|BUD|SBA|ALB|GRC|RMA|BUL|BUE|BUS|AGS|EMD|CON|ANK|SMY|ARM|SYR|BLA|SEV|MOS|UKR|STP|SPN|SPS|LVA|WAR|FIN|SWE|NOR|GOB|BAS|BAL|PRU|SIL|GAL|DEN|SKA|HEL|KIE|BER|MUN|RUH|NONE)$`)

func (d *Territory) validate(s string) error {
	if matched := reTerritory.MatchString(s); matched == false {
		return errors.New("Invalid value for Territory")
	}
	return nil
}

func (d *Territory) assign(s string) {
	*d = Territory(s)
}

func (d *Territory) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if err = d.validate(s); err == nil {
			d.assign(s)
		}
		return err
	}
	return err
}

func (d Territory) Value() (driver.Value, error) {
	return string(d), nil
}

func (z *Territory) Scan(s interface{}) (err error) {
	if z == nil {
		return errors.New("Territory: Scan on nil pointer")
	}
	if s == nil {
		return nil
	}
	*z = Territory(string(s.([]uint8)))
	return nil
}
