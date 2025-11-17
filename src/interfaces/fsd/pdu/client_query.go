// Package pdu
package pdu

import (
	"strings"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
)

type ClientQuery struct {
	*Base
	QueryType QueryType
	Payload   []string
}

func NewClientQuery(to string, queryType QueryType, payload ...string) *ClientQuery {
	return &ClientQuery{
		Base:      NewBase(fsd.ClientCommandClientQuery, global.FSDServerName, to),
		QueryType: queryType,
		Payload:   payload,
	}
}

func (c *ClientQuery) Build() []byte {
	return MakeProtocolDataUnitPacket(
		c.GetType(),
		c.From,
		c.To,
		c.QueryType.Data,
		strings.Join(c.Payload, Delimiter),
	)
}

func (c *ClientQuery) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZYSH_CTR SERVER FP  CPA421
	// [  0   ] [  1 ] [2] [  3 ]
	// ZYSH_CTR @94835 FA  CPA421 31100
	// [  0   ] [  1 ] [2] [  3 ] [ 4 ]
	// ZSHA_CTR @94835 BC  CXA8872 5074
	// [  0   ] [  1 ] [2] [  3  ] [ 4]
	if r := c.CheckLength(fsd.ClientCommandClientQuery, len(data)); r != nil {
		return nil, r
	}
	command := &ClientQuery{
		Base:      NewBase(fsd.ClientCommandClientQuery, data[0], data[1]),
		QueryType: QueryTypeInvalid,
		Payload:   data[3:],
	}
	if QueryTypes.IsValidEnum(data[2]) {
		command.QueryType = QueryTypes.GetEnum(data[2])
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
