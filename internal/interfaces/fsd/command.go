// Package fsd
package fsd

type ClientCommand string

var (
	AddAtc          = ClientCommand("#AA")
	RemoveAtc       = ClientCommand("#DA")
	AddPilot        = ClientCommand("#AP")
	RemovePilot     = ClientCommand("#DP")
	ProController   = ClientCommand("#PC")
	PilotPosition   = ClientCommand("@")
	AtcPosition     = ClientCommand("%")
	AtcSubVisPoint  = ClientCommand("'")
	Message         = ClientCommand("#TM")
	WeatherQuery    = ClientCommand("$AX")
	WeatherResponse = ClientCommand("$AR")
	SquawkBox       = ClientCommand("#SB")
	RequestHandoff  = ClientCommand("$HO")
	AcceptHandoff   = ClientCommand("$HA")
	Plan            = ClientCommand("$FP")
	AtcEditPlan     = ClientCommand("$AM")
	KillClient      = ClientCommand("$!!")
	Error           = ClientCommand("$ER")
	ClientQuery     = ClientCommand("$CQ")
	ClientResponse  = ClientCommand("$CR")
	TempData        = ClientCommand("$TD")
)

type CommandRequirement struct {
	RequireLength int
	Fatal         bool
}

func (c ClientCommand) String() string {
	return string(c)
}

func (c ClientCommand) Index() int {
	return 0
}
