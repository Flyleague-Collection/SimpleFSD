// Package entity
package entity

import "time"

type ActivityFacility struct {
	ID         uint                `gorm:"primarykey" json:"id"`
	ActivityId uint                `gorm:"index;not null" json:"activity_id"`
	MinRating  int                 `gorm:"default:2;not null" json:"min_rating"`
	Callsign   string              `gorm:"size:16;not null" json:"callsign"`
	Frequency  string              `gorm:"size:16;not null" json:"frequency"`
	Tier2Tower bool                `gorm:"default:false;not null" json:"tier2_tower"`
	SortIndex  int                 `gorm:"default:0;not null" json:"-"`
	Controller *ActivityController `gorm:"foreignKey:FacilityId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"controller"`
	CreatedAt  time.Time           `json:"-"`
	UpdatedAt  time.Time           `json:"-"`
}

func (facility *ActivityFacility) GetId() uint {
	return facility.ID
}

// Equal 比较两个ActivityFacility实例是否相等
// 通过逐一比较所有关键字段来判断两个设施是否相同
// 参数:
//   - other: 要与当前实例比较的另一个ActivityFacility实例
//
// 返回值:
//   - bool: 如果两个实例的所有字段都相等则返回true，否则返回false
func (facility *ActivityFacility) Equal(other *ActivityFacility) bool {
	if facility == nil || other == nil {
		return false
	}

	return facility.ID == other.ID &&
		facility.ActivityId == other.ActivityId &&
		facility.MinRating == other.MinRating &&
		facility.Callsign == other.Callsign &&
		facility.Frequency == other.Frequency &&
		facility.Tier2Tower == other.Tier2Tower &&
		facility.SortIndex == other.SortIndex
}

// Diff 比较两个ActivityFacility实例的差异
// 返回一个映射，其中包含当前实例与另一个实例之间不同字段的键值对
// 参数:
//   - other: 要与当前实例比较的另一个ActivityFacility实例
//
// 返回值:
//   - map[string]interface{}: 包含差异字段的映射，如果other为nil则返回nil
func (facility *ActivityFacility) Diff(other *ActivityFacility) map[string]interface{} {
	if facility == nil || other == nil {
		return nil
	}

	result := make(map[string]interface{})

	if other.MinRating >= 0 && facility.MinRating != other.MinRating {
		result["min_rating"] = other.MinRating
	}
	if other.Callsign != "" && facility.Callsign != other.Callsign {
		result["callsign"] = other.Callsign
	}
	if other.Frequency != "" && facility.Frequency != other.Frequency {
		result["frequency"] = other.Frequency
	}
	if facility.Tier2Tower != other.Tier2Tower {
		result["tier2_tower"] = other.Tier2Tower
	}
	if other.SortIndex >= 0 && facility.SortIndex != other.SortIndex {
		result["sort_index"] = other.SortIndex
	}

	return result
}
