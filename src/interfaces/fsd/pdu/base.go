// Package pdu Protocol Data Unit 协议数据单元
package pdu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type Interface interface {
	global.Builder[[]byte]
	GetType() fsd.ClientCommand
	Parse(data []string, raw []byte) (Interface, *fsd.CommandResult)
}

type Base struct {
	commandType fsd.ClientCommand
	raw         []byte
	From        string
	To          string
}

func NewBase(commandType fsd.ClientCommand, from string, to string) *Base {
	return &Base{
		commandType: commandType,
		From:        from,
		To:          to,
		raw:         nil,
	}
}

func (b *Base) GetType() fsd.ClientCommand {
	return b.commandType
}

func (b *Base) CheckLength(command fsd.ClientCommand, length int) *fsd.CommandResult {
	if length < command.Data.RequireLength {
		return fsd.CommandSyntaxError(command.Data.Fatal, "length", fmt.Errorf("command too short, expect %d but got %d", command.Data.RequireLength, length))
	}
	return nil
}

func (b *Base) Build() []byte {
	return b.raw
}

func (b *Base) Parse(_ []string, _ []byte) (Interface, *fsd.CommandResult) {
	return nil, fsd.CommandSyntaxError(true, "command", errors.New("this command only allow server to send to client"))
}

var CommandMap = map[fsd.ClientCommand]Interface{
	fsd.ClientCommandAddAtc:           &AddAtc{},
	fsd.ClientCommandRemoveAtc:        &RemoveAtc{},
	fsd.ClientCommandAddPilot:         &AddPilot{},
	fsd.ClientCommandRemovePilot:      &RemovePilot{},
	fsd.ClientCommandProController:    &ProController{},
	fsd.ClientCommandPilotPosition:    &PilotPosition{},
	fsd.ClientCommandAtcPosition:      &AtcPosition{},
	fsd.ClientCommandAtcSubVisPoint:   &AtcSubPosition{},
	fsd.ClientCommandMessage:          &Message{},
	fsd.ClientCommandWeatherQuery:     &WeatherQuery{},
	fsd.ClientCommandWeatherResponse:  &WeatherResponse{},
	fsd.ClientCommandPlaneInfo:        &PlaneInfo{},
	fsd.ClientCommandRequestHandoff:   &RequestHandoff{},
	fsd.ClientCommandAcceptHandoff:    &AcceptHandoff{},
	fsd.ClientCommandPlan:             &Plan{},
	fsd.ClientCommandAtcEditPlan:      &AtcEditPlan{},
	fsd.ClientCommandKillClient:       &KillClient{},
	fsd.ClientCommandError:            &Error{},
	fsd.ClientCommandClientQuery:      &ClientQuery{},
	fsd.ClientCommandClientResponse:   &ClientResponse{},
	fsd.ClientCommandClientIdent:      &ClientIdent{},
	fsd.ClientCommandServerIdent:      &ServerIdent{},
	fsd.ClientCommandSendFastPosition: &FastPosition{},
	fsd.ClientCommandFastPositionStop: &FastPositionStop{},
	fsd.ClientCommandFastPositionSlow: &FastPositionSlow{},
	fsd.ClientCommandFastPositionFast: &FastPositionFast{},
	fsd.ClientCommandUnknown:          nil,
}

const (
	PacketDelimiter = "\r\n"
	Delimiter       = ":"
)

func MakeProtocolDataUnitPacket(commandType fsd.ClientCommand, parts ...string) []byte {
	if len(parts) == 0 {
		return []byte(commandType.Value)
	}

	commands := make([]string, len(parts))
	commands[0] = commandType.Value + parts[0]
	copy(commands[1:], parts[1:])
	return []byte(strings.Join(commands, Delimiter))
}
