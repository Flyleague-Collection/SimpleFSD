// Package command
package command

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
)

type CommandContent struct {
	logger              log.LoggerInterface
	application         *interfaces.ApplicationContent
	isSimulatorServer   bool
	refuseOutRange      bool
	metarManager        interfaces.MetarManagerInterface
	clientManager       fsd.ClientManagerInterface
	messageQueue        queue.MessageQueueInterface
	userOperation       operation.UserOperationInterface
	flightPlanOperation operation.FlightPlanOperationInterface
	auditLogOperation   operation.AuditLogOperationInterface
}

func NewCommandContent(
	logger log.LoggerInterface,
	application *interfaces.ApplicationContent,
) *CommandContent {
	config := application.ConfigManager().Config()
	return &CommandContent{
		logger:              log.NewLoggerAdapter(logger, "CommandHandler"),
		application:         application,
		isSimulatorServer:   config.Server.General.SimulatorServer,
		refuseOutRange:      config.Server.FSDServer.RangeLimit.RefuseOutRange,
		metarManager:        application.MetarManager(),
		clientManager:       application.ClientManager(),
		messageQueue:        application.MessageQueue(),
		userOperation:       application.Operations().UserOperation(),
		flightPlanOperation: application.Operations().FlightPlanOperation(),
		auditLogOperation:   application.Operations().AuditLogOperation(),
	}
}
