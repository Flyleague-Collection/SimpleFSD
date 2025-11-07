// Package utils
package utils

import "math"

const (
	pitchMultiplier   = 256.0 / 90.0
	bankMultiplier    = 512.0 / 180.0
	headingMultiplier = 1024.0 / 360.0
)

func PackPBH(pitch, bank, heading float64, onGround bool) uint32 {
	pbh := uint32(0)

	if onGround {
		pbh |= 0b10
	}

	hdgVal := uint32(math.Round(heading*headingMultiplier)) & 0x3FF
	pbh |= hdgVal << 2

	bankVal := int(math.Round(bank * -bankMultiplier))
	pbh |= (uint32(bankVal) & 0x3FF) << 12

	pitchVal := int(math.Round(pitch * -pitchMultiplier))
	pbh |= (uint32(pitchVal) & 0x3FF) << 22

	return pbh
}

func UnpackPBH(pbh uint32) (pitch, bank, heading float64, onGround bool) {
	onGround = (pbh&0b10)>>1 == 1

	hdgBits := (pbh & 0xFFC) >> 2
	heading = float64(hdgBits) * (360.0 / 1024.0)

	bankVal := int32((pbh&0x3FF000)<<10) >> 22
	bank = float64(bankVal) * (-180.0 / 512.0)

	pitchVal := int32(pbh&0xFFC00000) >> 22
	pitch = float64(pitchVal) * (-90.0 / 256.0)

	return
}
