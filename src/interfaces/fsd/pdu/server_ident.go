// Package pdu
package pdu

type ServerIdent struct {
	*Base
	Name string
	Key  string
}

func (c *ServerIdent) Build() []byte {
	// SERVER CLIENT VATSIM FSD V3.53a 0815b2e12302
	// [  0 ] [  1 ] [        2      ] [     3    ]
	return []byte("$DISERVER:CLIENT:VATSIM FSD V3.53a:0815b2e12302")
}
