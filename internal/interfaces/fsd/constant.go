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

var PossibleClientCommands = [][]byte{[]byte(PilotPosition), []byte(AtcPosition), []byte(AtcSubVisPoint),
	[]byte(Message), []byte(ClientQuery), []byte(ClientResponse), []byte(Plan), []byte(AtcEditPlan), []byte(RequestHandoff),
	[]byte(AcceptHandoff), []byte(ProController), []byte(SquawkBox), []byte(AddAtc), []byte(RemoveAtc), []byte(AddPilot),
	[]byte(RemovePilot), []byte(KillClient)}

var CommandRequirements = map[ClientCommand]*CommandRequirement{
	AddAtc:         {12, true},
	AddPilot:       {8, true},
	AtcPosition:    {8, false},
	PilotPosition:  {10, false},
	AtcSubVisPoint: {4, false},
	ClientQuery:    {3, false},
	ClientResponse: {3, false},
	Message:        {3, false},
	Plan:           {17, false},
	AtcEditPlan:    {18, false},
	KillClient:     {2, false},
	RequestHandoff: {3, false},
	AcceptHandoff:  {3, false},
	ProController:  {3, false},
	SquawkBox:      {2, false},
}

var FacilityMap = map[string]Facility{
	"SUP": OBS,
	"OBS": OBS,
	"DEL": DEL,
	"GND": GND,
	"TWR": TWR,
	"APP": APP,
	"CTR": CTR,
	"FSS": FSS,
}
