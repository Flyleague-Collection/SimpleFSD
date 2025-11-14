// Package repository
package repository

import "github.com/half-nothing/simple-fsd/src/interfaces/database/entity"

type ControllerInterface interface {
	GetTotal() (total int64, err error)
	GetPage(pageNumber int, pageSize int) ([]*entity.User, int64, error)
	SetRating(user *entity.User, updateInfo map[string]interface{}) error
}
