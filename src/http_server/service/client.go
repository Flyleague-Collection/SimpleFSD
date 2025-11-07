// Package service
// 存放 ClientServiceInterface 的实现
package service

import (
	"errors"
	"fmt"

	"github.com/half-nothing/simple-fsd/src/interfaces"
	"github.com/half-nothing/simple-fsd/src/interfaces/config"
	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	. "github.com/half-nothing/simple-fsd/src/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/half-nothing/simple-fsd/src/interfaces/operation"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

type ClientService struct {
	logger            log.LoggerInterface
	clientManager     fsd.ClientManagerInterface
	messageQueue      queue.MessageQueueInterface
	config            *config.HttpServerConfig
	userOperation     operation.UserOperationInterface
	auditLogOperation operation.AuditLogOperationInterface
}

func NewClientService(
	logger log.LoggerInterface,
	config *config.HttpServerConfig,
	userOperation operation.UserOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
	clientManager fsd.ClientManagerInterface,
	messageQueue queue.MessageQueueInterface,
) *ClientService {
	service := &ClientService{
		logger:            log.NewLoggerAdapter(logger, "ClientService"),
		clientManager:     clientManager,
		config:            config,
		userOperation:     userOperation,
		auditLogOperation: auditLogOperation,
		messageQueue:      messageQueue,
	}
	return service
}

func (clientService *ClientService) GetOnlineClients() *fsd.OnlineClients {
	return clientService.clientManager.GetWhazzupContent()
}

func (clientService *ClientService) SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient] {
	if req.Uid <= 0 || req.SendTo == "" || req.Message == "" {
		return NewApiResponse[ResponseSendMessageToClient](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseSendMessageToClient](req.Permission, operation.ClientSendMessage); res != nil {
		return res
	}

	if err := clientService.messageQueue.SyncPublish(&queue.Message{
		Type: queue.SendMessageToClient,
		Data: &fsd.SendRawMessageData{
			From:    clientService.config.FormatCallsign(req.Cid),
			To:      req.SendTo,
			Message: req.Message,
		},
	}); err != nil {
		if errors.Is(err, fsd.ErrCallsignNotFound) {
			return NewApiResponse[ResponseSendMessageToClient](ErrClientNotFound, nil)
		}
		return NewApiResponse[ResponseSendMessageToClient](ErrSendMessage, nil)
	}

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: clientService.auditLogOperation.NewAuditLog(
			operation.ClientMessage,
			req.Cid,
			fmt.Sprintf("%s(%s)", req.SendTo, req.Message),
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	data := ResponseSendMessageToClient(true)
	return NewApiResponse[ResponseSendMessageToClient](SuccessSendMessage, &data)
}

func (clientService *ClientService) KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient] {
	if req.Uid <= 0 || req.TargetCallsign == "" {
		return NewApiResponse[ResponseKillClient](ErrIllegalParam, nil)
	}

	user, res := CheckPermissionFromDatabase[ResponseKillClient](clientService.userOperation, req.Uid, operation.ClientKill)
	if res != nil {
		return res
	}

	client, err := clientService.clientManager.KickClientFromServer(req.TargetCallsign, req.Reason)
	if err != nil {
		// KickClientFromServer目前仅返回ErrCallsignNotFound错误
		if errors.Is(err, fsd.ErrCallsignNotFound) {
			return NewApiResponse[ResponseKillClient](ErrClientNotFound, nil)
		}
		return NewApiResponse[ResponseKillClient](ErrUnknownServerError, nil)
	}

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.SendKickedFromServerEmail,
		Data: &interfaces.KickedFromServerEmailData{
			User:     client.User(),
			Operator: user,
			Reason:   req.Reason,
		},
	})

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: clientService.auditLogOperation.NewAuditLog(
			operation.ClientKicked,
			req.Cid,
			fmt.Sprintf("%s(%s)", req.TargetCallsign, req.Reason),
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	data := ResponseKillClient(true)
	return NewApiResponse[ResponseKillClient](SuccessKillClient, &data)
}

func (clientService *ClientService) GetClientFlightPath(req *RequestClientPath) *ApiResponse[ResponseClientPath] {
	if req.Callsign == "" {
		return NewApiResponse[ResponseClientPath](ErrIllegalParam, nil)
	}

	client, exist := clientService.clientManager.GetClient(req.Callsign)
	if !exist {
		return NewApiResponse[ResponseClientPath](ErrClientNotFound, nil)
	}

	data := ResponseClientPath(client.Paths())
	return NewApiResponse(SuccessGetClientPath, &data)
}

func (clientService *ClientService) SendBroadcastMessage(req *RequestSendBroadcastMessage) *ApiResponse[ResponseSendBroadcastMessage] {
	if req.Message == "" || !fsd.IsValidBroadcastTarget(req.Target) {
		return NewApiResponse[ResponseSendBroadcastMessage](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseSendBroadcastMessage](req.Permission, operation.ClientSendBroadcastMessage); res != nil {
		return res
	}

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.BroadcastMessage,
		Data: &fsd.BroadcastMessageData{
			From:    clientService.config.FormatCallsign(req.Cid),
			Target:  fsd.BroadcastTarget(req.Target),
			Message: req.Message,
		},
	})

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: clientService.auditLogOperation.NewAuditLog(
			operation.ClientBroadcastMessage,
			req.Cid,
			req.Target,
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: operation.ValueNotAvailable,
				NewValue: req.Message,
			},
		),
	})

	data := ResponseSendBroadcastMessage(true)
	return NewApiResponse[ResponseSendBroadcastMessage](SuccessSendBroadcastMessage, &data)
}
