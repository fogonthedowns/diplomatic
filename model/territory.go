package model

import "fmt"

const (
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
	NORWEGIAN_SEA             = Territory("NNS")
	MID_ATLANTIC_OCEAN        = Territory("MAO")
	BREST                     = Territory("BST")
	PICARDY                   = Territory("PIC")
	PARIS                     = Territory("PRS")
	BELGIUM                   = Territory("BGM")
	HOLLAND                   = Territory("HOL")
	GASCONY                   = Territory("GAS")
	BURGUNDY                  = Territory("BRG")
	MARSEILLES                = Territory("MAR")
	SPAIN                     = Territory("SPN")
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
)

type Territory string

// validNavyMovess defines a map of Valid moves for Navy Units
var validNavyMoves = map[Territory][]Territory{
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
	IONIAN_SEA:                []Territory{TYRRHENIAN_SEA, NAPLES, ADRIATIC_SEA, ALBANIA, GREECE, AEGEAN_SEA, EASTERN_MEDITERRANEAN},
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
	YORKSHIRE:                 []Territory{},
}

// validShipMovement will return true if the checked territory
// is included inside of the mapOfBorders map
// uses the origional terriotry as they key
func (t *Territory) validShipMovement(check Territory) bool {
	for _, borderTerritory := range validNavyMoves[*t] {
		if borderTerritory == check {
			return true
		}
	}
	return false
}
