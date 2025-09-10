// Package config
package config

import (
	"errors"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"strings"
)

var (
	ConfVersion, _ = newVersion(global.ConfigVersion)
	AppVersion, _  = newVersion(global.AppVersion)
)

func checkPort(port uint) *ValidResult {
	if port <= 0 {
		return ValidFail(errors.New("port must be greater than zero"))
	}
	if port > 65535 {
		return ValidFail(errors.New("port must be less than 65535"))
	}
	return ValidPass()
}

type checkVersionResult int

const (
	AllMatch checkVersionResult = iota
	MajorUnmatch
	MinorUnmatch
	PatchUnmatch
)

type Version struct {
	major   int
	minor   int
	patch   int
	version string
}

func newVersion(version string) (*Version, error) {
	versions := strings.Split(version, ".")
	if len(versions) < 3 {
		return nil, errors.New("invalid version String")
	}
	return &Version{
		major:   utils.StrToInt(versions[0], 0),
		minor:   utils.StrToInt(versions[1], 0),
		patch:   utils.StrToInt(versions[2], 0),
		version: version,
	}, nil
}

func (v *Version) checkVersion(version *Version) checkVersionResult {
	if v.major != version.major {
		return MajorUnmatch
	}
	if v.minor != version.minor {
		return MinorUnmatch
	}
	if v.patch != version.patch {
		return PatchUnmatch
	}
	return AllMatch
}

func (v *Version) String() string {
	return v.version
}
