// Package repository
package repository

import "github.com/half-nothing/simple-fsd/src/interfaces/database/entity"

type ControllerInterface interface {
	GetTotalControllers() (total int64, err error)
	GetControllers(page, pageSize int) (users []*entity.User, total int64, err error)
	SetControllerRating(user *entity.User, updateInfo map[string]interface{}) (err error)
}
