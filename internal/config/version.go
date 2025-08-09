// Package config
package config

import (
	"errors"
	"github.com/half-nothing/fsd-server/internal/utils"
	"strings"
)

type VersionType int

const (
	AllMatch VersionType = iota
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

func (v *Version) checkVersion(version *Version) VersionType {
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
