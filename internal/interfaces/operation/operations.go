// Package operation
package operation

type DatabaseOperations struct {
	userOperation       UserOperationInterface
	flightPlanOperation FlightPlanOperationInterface
}

func NewDatabaseOperations(
	userOperation UserOperationInterface,
	flightPlanOperation FlightPlanOperationInterface,
) *DatabaseOperations {
	return &DatabaseOperations{
		userOperation:       userOperation,
		flightPlanOperation: flightPlanOperation,
	}
}

func (db *DatabaseOperations) UserOperation() UserOperationInterface {
	return db.userOperation
}

func (db *DatabaseOperations) FlightPlanOperation() FlightPlanOperationInterface {
	return db.flightPlanOperation
}
