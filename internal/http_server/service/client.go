// Package service
// 存放 ClientServiceInterface 的实现
package service

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
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
		logger:            logger,
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

var (
	ErrSendMessage     = NewApiStatus("FAIL_SEND_MESSAGE", "发送消息失败", ServerInternalError)
	ErrClientNotFound  = NewApiStatus("CLIENT_NOT_FOUND", "指定客户端不存在", NotFound)
	SuccessSendMessage = NewApiStatus("SEND_MESSAGE", "发送成功", Ok)
)

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
			From:    req.Cid,
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

var SuccessKillClient = NewApiStatus("KILL_CLIENT", "成功踢出客户端", Ok)

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

	if clientService.config.Email.Template.EnableKickedFromServerEmail {
		clientService.messageQueue.Publish(&queue.Message{
			Type: queue.SendKickedFromServerEmail,
			Data: &KickedFromServerEmailData{
				User:     client.User(),
				Operator: user,
				Reason:   req.Reason,
			},
		})
	}

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

var SuccessGetClientPath = NewApiStatus("GET_CLIENT_PATH", "获取客户端飞行路径", Ok)

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
