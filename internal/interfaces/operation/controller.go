// Package operation
package operation

import "time"

type ControllerOperationInterface interface {
	GetTotalControllers() (total int64, err error)
	GetControllers(page, pageSize int) (users []*User, total int64, err error)
	SetControllerRating(user *User, rating int) (err error)
	SetControllerSolo(user *User, untilTime time.Time) (err error)
	UnsetControllerSolo(user *User) (err error)
	SetControllerUnderMonitor(user *User, underMonitor bool) (err error)
	SetControllerGuest(user *User, guest bool) (err error)
	SetControllerGuestRating(user *User, rating int) (err error)
}
