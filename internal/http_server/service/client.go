// Package service
package service

import (
	"errors"
	"fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"time"
)

type ClientService struct {
	logger            log.LoggerInterface
	onlineClient      *utils.CachedValue[OnlineClients]
	clientManager     fsd.ClientManagerInterface
	emailService      EmailServiceInterface
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
	emailService EmailServiceInterface,
) *ClientService {
	service := &ClientService{
		logger:            logger,
		clientManager:     clientManager,
		emailService:      emailService,
		config:            config,
		userOperation:     userOperation,
		auditLogOperation: auditLogOperation,
	}
	service.onlineClient = utils.NewCachedValue[OnlineClients](config.CacheDuration, func() *OnlineClients { return service.getOnlineClient() })
	return service
}

func (clientService *ClientService) getOnlineClient() *OnlineClients {
	data := &OnlineClients{
		General: OnlineGeneral{
			Version:          3,
			ConnectedClients: 0,
			OnlinePilot:      0,
			OnlineController: 0,
		},
		Pilots:      make([]*OnlinePilot, 0),
		Controllers: make([]*OnlineController, 0),
	}

	clientCopy := clientService.clientManager.GetClientSnapshot()
	defer clientService.clientManager.PutSlice(clientCopy)

	for _, client := range clientCopy {
		if client == nil || client.Disconnected() {
			continue
		}
		data.General.ConnectedClients++
		if client.IsAtc() {
			data.General.OnlineController++
			controller := &OnlineController{
				Cid:         client.User().Cid,
				Callsign:    client.Callsign(),
				RealName:    client.RealName(),
				Latitude:    client.Position()[0].Latitude,
				Longitude:   client.Position()[0].Longitude,
				Rating:      client.Rating().Index(),
				Facility:    client.Facility().Index(),
				Frequency:   client.Frequency() + 100000,
				Range:       int(client.VisualRange()),
				OfflineTime: client.LogoffTime(),
				IsBreak:     client.IsBreak(),
				AtcInfo:     client.AtisInfo(),
				LogonTime:   client.History().StartTime.Format(time.DateTime),
			}
			data.Controllers = append(data.Controllers, controller)
		} else {
			data.General.OnlinePilot++
			pilot := &OnlinePilot{
				Cid:         client.User().Cid,
				Callsign:    client.Callsign(),
				RealName:    client.RealName(),
				Latitude:    client.Position()[0].Latitude,
				Longitude:   client.Position()[0].Longitude,
				Transponder: client.Transponder(),
				Heading:     client.Heading(),
				Altitude:    client.Altitude(),
				GroundSpeed: client.GroundSpeed(),
				FlightPlan:  client.FlightPlan(),
				LogonTime:   client.History().StartTime.Format(time.DateTime),
			}
			data.Pilots = append(data.Pilots, pilot)
		}
	}

	data.General.GenerateTime = time.Now().Format(time.DateTime)

	return data
}

func (clientService *ClientService) GetOnlineClient() *OnlineClients {
	return clientService.onlineClient.GetValue()
}

var (
	ErrSendMessage      = ApiStatus{StatusName: "FAIL_SEND_MESSAGE", Description: "发送消息失败", HttpCode: ServerInternalError}
	ErrCallsignNotFound = ApiStatus{StatusName: "CALLSIGN_NOT_FOUND", Description: "发送目标不在线", HttpCode: NotFound}
	SuccessSendMessage  = ApiStatus{StatusName: "SEND_MESSAGE", Description: "发送成功", HttpCode: Ok}
)

func (clientService *ClientService) SendMessageToClient(req *RequestSendMessageToClient) *ApiResponse[ResponseSendMessageToClient] {
	if req.Uid <= 0 || req.SendTo == "" || req.Message == "" {
		return NewApiResponse[ResponseSendMessageToClient](&ErrIllegalParam, Unsatisfied, nil)
	}
	if req.Permission <= 0 {
		return NewApiResponse[ResponseSendMessageToClient](&ErrNoPermission, Unsatisfied, nil)
	}
	permission := operation.Permission(req.Permission)
	if !permission.HasPermission(operation.ClientSendMessage) {
		return NewApiResponse[ResponseSendMessageToClient](&ErrNoPermission, Unsatisfied, nil)
	}
	if err := clientService.clientManager.SendRawMessageTo(req.Cid, req.SendTo, req.Message); err != nil {
		if errors.Is(err, fsd.ErrCallsignNotFound) {
			return NewApiResponse[ResponseSendMessageToClient](&ErrCallsignNotFound, Unsatisfied, nil)
		}
		return NewApiResponse[ResponseSendMessageToClient](&ErrSendMessage, Unsatisfied, nil)
	}

	go func() {
		auditLog := clientService.auditLogOperation.NewAuditLog(operation.ClientMessage, req.Cid,
			fmt.Sprintf("%s(%s)", req.SendTo, req.Message), req.Ip, req.UserAgent, nil)
		err := clientService.auditLogOperation.SaveAuditLog(auditLog)
		if err != nil {
			clientService.logger.ErrorF("Fail to create audit log for client_message, detail: %v", err)
		}
	}()

	data := ResponseSendMessageToClient(true)
	return NewApiResponse[ResponseSendMessageToClient](&SuccessSendMessage, Unsatisfied, &data)
}

var SuccessKillClient = ApiStatus{StatusName: "KILL_CLIENT", Description: "成功踢出客户端", HttpCode: Ok}

func (clientService *ClientService) KillClient(req *RequestKillClient) *ApiResponse[ResponseKillClient] {
	if req.Uid <= 0 || req.TargetCallsign == "" {
		return NewApiResponse[ResponseKillClient](&ErrIllegalParam, Unsatisfied, nil)
	}
	user, res := CallDBFuncAndCheckError[operation.User, ResponseKillClient](func() (*operation.User, error) {
		return clientService.userOperation.GetUserByUid(req.Uid)
	})
	if res != nil {
		return res
	}
	permission := operation.Permission(user.Permission)
	if !permission.HasPermission(operation.ClientKill) {
		return NewApiResponse[ResponseKillClient](&ErrNoPermission, Unsatisfied, nil)
	}
	client, ok := clientService.clientManager.GetClient(req.TargetCallsign)
	if !ok {
		return NewApiResponse[ResponseKillClient](&ErrCallsignNotFound, Unsatisfied, nil)
	}
	client.MarkedDisconnect(false)

	go func() {
		if clientService.config.Email.Template.EnableKickedFromServerEmail {
			if err := clientService.emailService.SendKickedFromServerEmail(client.User(), user, req.Reason); err != nil {
				clientService.logger.ErrorF("SendRatingChangeEmail Failed: %v", err)
			}
		}
	}()

	go func() {
		auditLog := clientService.auditLogOperation.NewAuditLog(operation.ClientKicked, req.Cid,
			fmt.Sprintf("%s(%s)", req.TargetCallsign, req.Reason), req.Ip, req.UserAgent, nil)
		err := clientService.auditLogOperation.SaveAuditLog(auditLog)
		if err != nil {
			clientService.logger.ErrorF("Fail to create audit log for client_kicked, detail: %v", err)
		}
	}()

	data := ResponseKillClient(true)
	return NewApiResponse[ResponseKillClient](&SuccessKillClient, Unsatisfied, &data)
}

var (
	ErrClientNotFound    = ApiStatus{StatusName: "CLIENT_NOT_FOUND", Description: "指定客户端不存在", HttpCode: NotFound}
	SuccessGetClientPath = ApiStatus{StatusName: "GET_CLIENT_PATH", Description: "获取指定客户端飞行路径", HttpCode: Ok}
)

func (clientService *ClientService) GetClientPath(req *RequestClientPath) *ApiResponse[ResponseClientPath] {
	if req.Callsign == "" {
		return NewApiResponse[ResponseClientPath](&ErrIllegalParam, Unsatisfied, nil)
	}
	client, exist := clientService.clientManager.GetClient(req.Callsign)
	if !exist {
		return NewApiResponse[ResponseClientPath](&ErrClientNotFound, Unsatisfied, nil)
	}
	data := ResponseClientPath(client.Paths())
	return NewApiResponse(&SuccessGetClientPath, Unsatisfied, &data)
}
