// Package fsd
package fsd

import "strings"

type RatingModel struct {
	Id        int    `json:"id"`
	ShortName string `json:"short_name"`
	LongName  string `json:"long_name"`
}
type Rating int

const (
	Ban Rating = iota - 1
	Normal
	Observer
	STU1
	STU2
	STU3
	CTR1
	CTR2
	CTR3
	Instructor1
	Instructor2
	Instructor3
	Supervisor
	Administrator
)

var Ratings = []*RatingModel{
	{-1, "Baned", "Suspended"},
	{0, "Normal", "Normal"},
	{1, "OBS", "Observer"},
	{2, "S1", "Tower Trainee"},
	{3, "S2", "Tower Controller"},
	{4, "S3", "Senior Student"},
	{5, "C1", "Enroute Controller"},
	{6, "C2", "Controller 2 (not in use)"},
	{7, "C3", "Senior Controller"},
	{8, "I1", "Instructor"},
	{9, "I2", "Instructor 2 (not in use)"},
	{10, "I3", "Senior Instructor"},
	{11, "SUP", "Supervisor"},
	{12, "ADM", "Administrator"},
}

func (r Rating) String() string {
	return Ratings[r+1].ShortName
}

func (r Rating) Index() int {
	return int(r)
}

func ToRatingString(rating int, tier2 bool, um bool, solo bool) string {
	result := []string{Rating(rating).String()}
	if tier2 {
		result = append(result, "Tier2")
	}
	if um {
		result = append(result, "UM")
	}
	if solo {
		result = append(result, "SOLO")
	}
	return strings.Join(result, " ")
}

func IsValidRating(r int) bool {
	return int(Ban) <= r && r <= int(Administrator)
}
