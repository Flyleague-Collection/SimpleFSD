// Package pdu
package pdu

import (
	"strings"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type ClientResponse struct {
	*Base
	QueryType QueryType
	Payload   []string
}

func NewClientResponse(to string, queryType QueryType, payload ...string) *ClientResponse {
	return &ClientResponse{
		Base:      NewBase(fsd.ClientCommandClientResponse, global.FSDServerName, to),
		QueryType: queryType,
		Payload:   payload,
	}
}

func (c ClientResponse) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		c.From,
		c.To,
		c.QueryType.Data,
		strings.Join(c.Payload, Delimiter),
	)
}

func (c ClientResponse) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR ZSSS_APP CAPS ATCINFO=1 SECPOS=1 MODELDESC=1 ONGOINGCOORD=1 NEWINFO=1 TEAMSPEAK=1 ICAOEQ=1
	// [   0  ] [   1  ] [ 2] [   3   ] [  4   ] [    5    ] [     6      ] [   7   ] [     8   ] [  9   ]
	// ZSHA_CTR SERVER ATIS  T  ZSHA_CTR Shanghai Control
	// [   0  ] [  1 ] [ 2] [3] [           4           ]
	if r := c.CheckLength(fsd.ClientCommandClientResponse, len(data)); r != nil {
		return nil, r
	}
	command := &ClientResponse{
		Base:      NewBase(fsd.ClientCommandClientResponse, data[0], data[1]),
		QueryType: QueryTypeInvalid,
		Payload:   data[3:],
	}
	if QueryTypes.IsValidEnum(data[2]) {
		command.QueryType = QueryTypes.GetEnum(data[2])
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
