// Package utils
package utils

import (
	"math"
	"testing"
)

func IsEqual(f1, f2 float64) bool {
	if f1 > f2 {
		return math.Dim(f1, f2) < 0.000001
	} else {
		return math.Dim(f2, f1) < 0.000001
	}
}

func TestPackPBH(t *testing.T) {
	tests := []struct {
		pitch    float64
		bank     float64
		heading  float64
		onGround bool
		expected uint32
	}{
		{0.35156250, 0.35156250, 352.26562500, true, 4294967210},
		{0.35156250, 0, 172.26562500, true, 4290774954},
		{20.03906250, 0, 176.13281250, false, 4055893972},
		{55.19531250, 116.36718750, 298.47656250, false, 3639303492},
	}

	pass := 0
	fail := 0
	for _, test := range tests {
		pbh := PackPBH(test.pitch, test.bank, test.heading, test.onGround)
		if pbh != test.expected {
			fail++
			t.Errorf("PackPBH(%f, %f, %f, %v) = %d; expected %d", test.pitch, test.bank, test.heading, test.onGround, pbh, test.expected)
			continue
		}
		pass++
	}
	t.Logf("TestPackPBH: %d pass, %d fail", pass, fail)

}

func TestUnpackPBH(t *testing.T) {
	tests := []struct {
		input    uint32
		pitch    float64
		bank     float64
		heading  float64
		onGround bool
	}{
		{4294967210, 0.35156250, 0.35156250, 352.26562500, true},
		{4290774954, 0.35156250, 0, 172.26562500, true},
		{4055893972, 20.03906250, 0, 176.13281250, false},
		{3639303492, 55.19531250, 116.36718750, 298.47656250, false},
	}
	pass := 0
	fail := 0
	for _, test := range tests {
		pitch, bank, heading, onGround := UnpackPBH(test.input)
		if !IsEqual(pitch, test.pitch) || !IsEqual(bank, test.bank) || !IsEqual(heading, test.heading) || onGround != test.onGround {
			fail++
			t.Errorf("UnpackPBH(%d) = %.8f, %.8f, %.8f, %v; expected %.8f, %.8f, %.8f, %v", test.input, pitch, bank, heading, onGround, test.pitch, test.bank, test.heading, test.onGround)
			continue
		}
		pass++
	}
	t.Logf("TestUnpackPBH: %d pass, %d fail", pass, fail)
}
