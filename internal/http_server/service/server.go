// Package service
package service

import (
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/half-nothing/simple-fsd/internal/utils"
)

type ServerService struct {
	logger            log.LoggerInterface
	config            *config.ServerConfig
	userOperation     operation.UserOperationInterface
	activityOperation operation.ActivityOperationInterface
	serverConfig      *utils.CachedValue[ResponseGetServerConfig]
	serverInfo        *utils.CachedValue[ResponseGetServerInfo]
	serverOnlineTime  *utils.CachedValue[ResponseGetTimeRating]
}

func NewServerService(
	logger log.LoggerInterface,
	config *config.ServerConfig,
	userOperation operation.UserOperationInterface,
	activityOperation operation.ActivityOperationInterface,
) *ServerService {
	service := &ServerService{
		logger:            logger,
		config:            config,
		userOperation:     userOperation,
		activityOperation: activityOperation,
	}
	service.serverConfig = utils.NewCachedValue[ResponseGetServerConfig](0, func() *ResponseGetServerConfig { return service.getServerConfig() })
	service.serverInfo = utils.NewCachedValue[ResponseGetServerInfo](config.HttpServer.CacheDuration, func() *ResponseGetServerInfo { return service.getServerInfo() })
	service.serverOnlineTime = utils.NewCachedValue[ResponseGetTimeRating](config.HttpServer.CacheDuration, func() *ResponseGetTimeRating { return service.getTimeRating() })
	return service
}

func (serverService *ServerService) getServerConfig() *ResponseGetServerConfig {
	return &ResponseGetServerConfig{
		Limits: &ServerLimits{
			UsernameLengthMin: serverService.config.HttpServer.Limits.UsernameLengthMin,
			UsernameLengthMax: serverService.config.HttpServer.Limits.UsernameLengthMax,
			PasswordLengthMin: serverService.config.HttpServer.Limits.PasswordLengthMin,
			PasswordLengthMax: serverService.config.HttpServer.Limits.PasswordLengthMax,
			EmailLengthMin:    serverService.config.HttpServer.Limits.EmailLengthMin,
			EmailLengthMax:    serverService.config.HttpServer.Limits.EmailLengthMax,
			CidMin:            serverService.config.HttpServer.Limits.CidMin,
			CidMax:            serverService.config.HttpServer.Limits.CidMax,
			SimulatorServer:   serverService.config.General.SimulatorServer,
		},
		ImageLimit: &ImageLimit{
			MaxAllowSize: int(serverService.config.HttpServer.Store.FileLimit.ImageLimit.MaxFileSize),
			AllowedExt:   serverService.config.HttpServer.Store.FileLimit.ImageLimit.AllowedFileExt,
		},
		Facilities: &fsd.Facilities,
		Ratings:    &fsd.Ratings,
	}
}

func (serverService *ServerService) getServerInfo() *ResponseGetServerInfo {
	totalUser, err := serverService.userOperation.GetTotalUsers()
	if err != nil {
		serverService.logger.ErrorF("ServerService.GetTotalUsers error: %v", err)
		totalUser = 0
	}
	totalControllers, err := serverService.userOperation.GetTotalControllers()
	if err != nil {
		serverService.logger.ErrorF("ServerService.GetTotalControllers error: %v", err)
		totalControllers = 0
	}
	totalActivities, err := serverService.activityOperation.GetTotalActivities()
	if err != nil {
		serverService.logger.ErrorF("ServerService.GetTotalActivities error: %v", err)
		totalActivities = 0
	}
	return &ResponseGetServerInfo{
		TotalUser:       totalUser,
		TotalController: totalControllers,
		TotalActivity:   totalActivities,
	}
}

func (serverService *ServerService) getTimeRating() *ResponseGetTimeRating {
	pilots, controllers, err := serverService.userOperation.GetTimeRatings()
	if err != nil {
		serverService.logger.ErrorF("ServerService.GetTimeRatings error: %v", err)
		return &ResponseGetTimeRating{}
	}
	data := &ResponseGetTimeRating{
		Pilots:      make([]OnlineTime, 0),
		Controllers: make([]OnlineTime, 0),
	}
	for _, pilot := range pilots {
		data.Pilots = append(data.Pilots, OnlineTime{
			Cid:  pilot.Cid,
			Time: pilot.TotalPilotTime,
		})
	}
	for _, controller := range controllers {
		data.Controllers = append(data.Controllers, OnlineTime{
			Cid:  controller.Cid,
			Time: controller.TotalAtcTime,
		})
	}
	return data
}

var SuccessGetServerConfig = ApiStatus{StatusName: "GET_SERVER_CONFIG", Description: "成功获取服务器配置", HttpCode: Ok}

func (serverService *ServerService) GetServerConfig() *ApiResponse[ResponseGetServerConfig] {
	return NewApiResponse(&SuccessGetServerConfig, Unsatisfied, serverService.serverConfig.GetValue())
}

var SuccessGetServerInfo = ApiStatus{StatusName: "GET_SERVER_INFO", Description: "成功获取服务器信息", HttpCode: Ok}

func (serverService *ServerService) GetServerInfo() *ApiResponse[ResponseGetServerInfo] {
	return NewApiResponse(&SuccessGetServerInfo, Unsatisfied, serverService.serverInfo.GetValue())
}

var SuccessGetTimeRating = ApiStatus{StatusName: "GET_TIME_RATING", Description: "成功获取服务器排行榜", HttpCode: Ok}

func (serverService *ServerService) GetTimeRating() *ApiResponse[ResponseGetTimeRating] {
	return NewApiResponse(&SuccessGetTimeRating, Unsatisfied, serverService.serverOnlineTime.GetValue())
}
