package fsd

type Enum interface {
	String() string
	Index() int
}

const EuroscopeFrequency = "@94835"

const AllowAtcFacility = DEL | GND | TWR | APP | CTR | FSS

var AllowKillRating = []Rating{Supervisor, Administrator}

var RatingFacilityMap = map[Rating]Facility{
	Ban:           0,
	Normal:        Pilot,
	Observer:      Pilot | OBS,
	STU1:          Pilot | OBS | DEL | GND,
	STU2:          Pilot | OBS | DEL | GND | TWR,
	STU3:          Pilot | OBS | DEL | GND | TWR | APP,
	CTR1:          Pilot | OBS | DEL | GND | TWR | APP | CTR,
	CTR2:          Pilot | OBS | DEL | GND | TWR | APP | CTR,
	CTR3:          Pilot | OBS | DEL | GND | TWR | APP | CTR | FSS,
	Instructor1:   Pilot | OBS | DEL | GND | TWR | APP | CTR | FSS,
	Instructor2:   Pilot | OBS | DEL | GND | TWR | APP | CTR | FSS,
	Instructor3:   Pilot | OBS | DEL | GND | TWR | APP | CTR | FSS,
	Supervisor:    Pilot | OBS | DEL | GND | TWR | APP | CTR | FSS,
	Administrator: Pilot | OBS | DEL | GND | TWR | APP | CTR | FSS,
}

var FacilityMap = map[string]Facility{
	"SUP":  OBS,
	"OBS":  OBS,
	"DEL":  DEL,
	"RMP":  GND,
	"GND":  GND,
	"TWR":  TWR,
	"APP":  APP,
	"CTR":  CTR,
	"FSS":  FSS,
	"ATIS": TWR,
}
