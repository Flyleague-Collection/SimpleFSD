// Package pdu
package pdu

import (
	"fmt"
	"strings"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd/rating"
	"github.com/half-nothing/simple-fsd/src/interfaces/global"
	"github.com/half-nothing/simple-fsd/src/utils"
)

type AtcPosition struct {
	*Base
	Frequencies     []int
	Facility        rating.Facility
	VisibilityRange int
	Rating          *rating.Rating
	Position        *global.Position
}

func (c *AtcPosition) Parse(data []string, raw []byte) (Interface, *fsd.CommandResult) {
	// ZSHA_CTR 24550  6  600  5  27.28025 118.28701  0
	// [   0  ] [ 1 ] [2] [3] [4] [   5  ] [   6   ] [7]
	if r := c.CheckLength(fsd.ClientCommandAtcPosition, len(data)); r != nil {
		return nil, r
	}
	command := &AtcPosition{
		Base: NewBase(fsd.ClientCommandAtcPosition, data[0], ""),
	}
	frequencies := make([]int, 0)
	utils.ForEach(strings.Split(data[1], "&"), func(index int, frequency string) {
		frequencies = append(frequencies, utils.StrToInt(frequency, -1)+100000)
	})
	if utils.Any(frequencies, func(frequency int) bool { return !fsd.FrequencyValid(frequency) }) {
		return nil, fsd.CommandSyntaxError(false, "frequency", fmt.Errorf("frequency %s is not vaild", data[1]))
	}
	command.Frequencies = frequencies
	facilityIndex := utils.StrToInt(data[2], -1)
	if facilityIndex == -1 {
		return nil, fsd.CommandSyntaxError(false, "facility", fmt.Errorf("facility %s is not a number", data[2]))
	}
	command.Facility = rating.Facility(1 << facilityIndex)
	if rating.Facilities[command.Facility] == nil {
		return nil, fsd.CommandSyntaxError(false, "facility", fmt.Errorf("facility %s not vaild", data[2]))
	}
	command.VisibilityRange = utils.StrToInt(data[3], -1)
	if command.VisibilityRange == -1 {
		return nil, fsd.CommandSyntaxError(false, "visibilityRange", fmt.Errorf("visibilityRange %s is not a number", data[3]))
	}
	reqRating := utils.StrToInt(data[4], -1)
	if reqRating == -1 || !rating.Ratings.IsValidEnum(reqRating) {
		return nil, fsd.CommandSyntaxError(false, "Rating", fmt.Errorf("rating(%s) not vaild", data[4]))
	}
	command.Rating = (*rating.Rating)(rating.Ratings.GetEnum(reqRating))
	command.Position = &global.Position{
		Latitude:  utils.StrToFloat(data[5], 0),
		Longitude: utils.StrToFloat(data[6], 0),
	}
	if !command.Position.Valid() {
		return nil, fsd.CommandSyntaxError(false, "position", fmt.Errorf("position %s %s is not a valid", data[5], data[6]))
	}
	command.raw = raw
	return command, fsd.CommandResultSuccess()
}
