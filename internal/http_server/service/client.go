// Package service
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

func (clientService *ClientService) GetOnlineClient() *fsd.OnlineClients {
	return clientService.clientManager.GetWhazzupContent()
}

var (
	ErrSendMessage      = ApiStatus{StatusName: "FAIL_SEND_MESSAGE", Description: "发送消息失败", HttpCode: ServerInternalError}
	ErrCallsignNotFound = ApiStatus{StatusName: "CALLSIGN_NOT_FOUND", Description: "发送目标不在线", HttpCode: NotFound}
	SuccessSendMessage  = ApiStatus{StatusName: "SEND_MESSAGE", Description: "发送成功", HttpCode: Ok}
)

func (clientService *ClientService) SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient] {
	if req.Uid <= 0 || req.SendTo == "" || req.Message == "" {
		return NewApiResponse[ResponseSendMessageToClient](ErrIllegalParam, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseSendMessageToClient](ErrNoPermission, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ClientSendMessage) {
		return NewApiResponse[ResponseSendMessageToClient](ErrNoPermission, nil)
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
			return NewApiResponse[ResponseSendMessageToClient](&ErrCallsignNotFound, nil)
		}
		return NewApiResponse[ResponseSendMessageToClient](&ErrSendMessage, nil)
	}

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: clientService.auditLogOperation.NewAuditLog(operation.ClientMessage, req.Cid,
			fmt.Sprintf("%s(%s)", req.SendTo, req.Message), req.Ip, req.UserAgent, nil),
	})

	data := ResponseSendMessageToClient(true)
	return NewApiResponse[ResponseSendMessageToClient](&SuccessSendMessage, &data)
}

var SuccessKillClient = ApiStatus{StatusName: "KILL_CLIENT", Description: "成功踢出客户端", HttpCode: Ok}

func (clientService *ClientService) KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient] {
	if req.Uid <= 0 || req.TargetCallsign == "" {
		return NewApiResponse[ResponseKillClient](ErrIllegalParam, nil)
	}
	user, res := CallDBFunc[operation.User, ResponseKillClient](func() (*operation.User, error) {
		return clientService.userOperation.GetUserByUid(req.Uid)
	})
	if res != nil {
		return res
	}
	permission := operation.Permission(user.Permission)
	if !permission.HasPermission(operation.ClientKill) {
		return NewApiResponse[ResponseKillClient](ErrNoPermission, nil)
	}
	client, ok := clientService.clientManager.GetClient(req.TargetCallsign)
	if !ok {
		return NewApiResponse[ResponseKillClient](&ErrCallsignNotFound, nil)
	}
	client.MarkedDisconnect(false)

	if clientService.config.Email.Template.EnableKickedFromServerEmail {
		clientService.messageQueue.Publish(&queue.Message{
			Type: queue.SendKickedFromServerEmail,
			Data: &SendKickedFromServerData{
				User:     client.User(),
				Operator: user,
				Reason:   req.Reason,
			},
		})
	}

	clientService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: clientService.auditLogOperation.NewAuditLog(operation.ClientKicked, req.Cid,
			fmt.Sprintf("%s(%s)", req.TargetCallsign, req.Reason), req.Ip, req.UserAgent, nil),
	})

	data := ResponseKillClient(true)
	return NewApiResponse[ResponseKillClient](&SuccessKillClient, &data)
}

var (
	ErrClientNotFound    = ApiStatus{StatusName: "CLIENT_NOT_FOUND", Description: "指定客户端不存在", HttpCode: NotFound}
	SuccessGetClientPath = ApiStatus{StatusName: "GET_CLIENT_PATH", Description: "获取指定客户端飞行路径", HttpCode: Ok}
)

func (clientService *ClientService) GetClientPath(req *RequestClientPath) *ApiResponse[ResponseClientPath] {
	if req.Callsign == "" {
		return NewApiResponse[ResponseClientPath](ErrIllegalParam, nil)
	}
	client, exist := clientService.clientManager.GetClient(req.Callsign)
	if !exist {
		return NewApiResponse[ResponseClientPath](&ErrClientNotFound, nil)
	}
	data := ResponseClientPath(client.Paths())
	return NewApiResponse(&SuccessGetClientPath, &data)
}
