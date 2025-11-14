// Package repository
package repository

import (
	"strings"
	"time"

	"github.com/half-nothing/simple-fsd/src/interfaces/database/entity"
)

type ActivityBuilder struct {
	activity *entity.Activity
}

func NewActivityBuilder() *ActivityBuilder {
	return &ActivityBuilder{
		activity: &entity.Activity{
			Status: ActivityStatusRegistering.Value,
		},
	}
}

// 活动通用信息

func (builder *ActivityBuilder) SetType(activityType ActivityType) *ActivityBuilder {
	builder.activity.Type = activityType.Value
	return builder
}

func (builder *ActivityBuilder) SetPublisher(publisherId int) *ActivityBuilder {
	builder.activity.Publisher = publisherId
	return builder
}

func (builder *ActivityBuilder) SetTitle(title string) *ActivityBuilder {
	builder.activity.Title = title
	return builder
}

func (builder *ActivityBuilder) SetImage(image *entity.Image) *ActivityBuilder {
	builder.activity.ImageId = image.ID
	return builder
}

func (builder *ActivityBuilder) SetActiveTime(activeTime time.Time) *ActivityBuilder {
	builder.activity.ActiveTime = activeTime
	return builder
}

func (builder *ActivityBuilder) SetNOTAMS(notams string) *ActivityBuilder {
	builder.activity.NOTAMS = notams
	return builder
}

// 单向单站活动信息

func (builder *ActivityBuilder) SetDepartureAirport(dep string) *ActivityBuilder {
	builder.activity.DepartureAirport = dep
	return builder
}

func (builder *ActivityBuilder) SetArrivalAirport(arr string) *ActivityBuilder {
	builder.activity.ArrivalAirport = arr
	return builder
}

func (builder *ActivityBuilder) SetRoute(route string) *ActivityBuilder {
	builder.activity.Route = route
	return builder
}

func (builder *ActivityBuilder) SetDistance(distance int) *ActivityBuilder {
	builder.activity.Distance = distance
	return builder
}

// 双向双站活动信息

func (builder *ActivityBuilder) SetRoute2(route string) *ActivityBuilder {
	builder.activity.Route2 = route
	return builder
}

func (builder *ActivityBuilder) SetDistance2(distance int) *ActivityBuilder {
	builder.activity.Distance2 = distance
	return builder
}

// 空域开放日

func (builder *ActivityBuilder) SetOpenFir(firs []string) *ActivityBuilder {
	builder.activity.DepartureAirport = strings.Join(firs, ",")
	return builder
}

func (builder *ActivityBuilder) Build() *entity.Activity {
	return builder.activity
}
