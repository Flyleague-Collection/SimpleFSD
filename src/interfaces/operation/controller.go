// Package operation
package operation

type ControllerOperationInterface interface {
	GetTotalControllers() (total int64, err error)
	GetControllers(page, pageSize int) (users []*User, total int64, err error)
	SetControllerRating(user *User, updateInfo map[string]interface{}) (err error)
}
