// Package service
// 存放 FlightPlanServiceInterface 的实现
package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/half-nothing/simple-fsd/src/interfaces/fsd"
	. "github.com/half-nothing/simple-fsd/src/interfaces/http/service"
	"github.com/half-nothing/simple-fsd/src/interfaces/log"
	"github.com/half-nothing/simple-fsd/src/interfaces/operation"
	"github.com/half-nothing/simple-fsd/src/interfaces/queue"
)

type FlightPlanService struct {
	logger              log.LoggerInterface
	messageQueue        queue.MessageQueueInterface
	userOperation       operation.UserOperationInterface
	flightPlanOperation operation.FlightPlanOperationInterface
	auditLogOperation   operation.AuditLogOperationInterface
}

func NewFlightPlanService(
	logger log.LoggerInterface,
	messageQueue queue.MessageQueueInterface,
	userOperation operation.UserOperationInterface,
	flightPlanOperation operation.FlightPlanOperationInterface,
	auditLogOperation operation.AuditLogOperationInterface,
) *FlightPlanService {
	return &FlightPlanService{
		logger:              log.NewLoggerAdapter(logger, "FlightPlanService"),
		messageQueue:        messageQueue,
		userOperation:       userOperation,
		flightPlanOperation: flightPlanOperation,
		auditLogOperation:   auditLogOperation,
	}
}

func (flightPlanService *FlightPlanService) SubmitFlightPlan(req *RequestSubmitFlightPlan) *ApiResponse[ResponseSubmitFlightPlan] {
	if req.FlightPlan == nil {
		return NewApiResponse[ResponseSubmitFlightPlan](ErrIllegalParam, nil)
	}

	if flightPlan, err := flightPlanService.flightPlanOperation.GetFlightPlanByCid(req.JwtHeader.Cid); err != nil {
		if errors.Is(err, operation.ErrFlightPlanNotFound) {
			req.FlightPlan.ID = 0
		} else {
			return NewApiResponse[ResponseSubmitFlightPlan](ErrDatabaseFail, nil)
		}
	} else {
		if flightPlan.Locked && flightPlan.DepartureAirport == req.DepartureAirport && flightPlan.ArrivalAirport == req.ArrivalAirport {
			return NewApiResponse[ResponseSubmitFlightPlan](ErrFlightPlanLocked, nil)
		}
		req.FlightPlan.Locked = false
		req.FlightPlan.ID = flightPlan.ID
		req.FlightPlan.CreatedAt = flightPlan.CreatedAt
	}

	req.FlightPlan.Cid = req.JwtHeader.Cid
	req.FlightPlan.FromWeb = true

	if res := CallDBFuncWithoutRet[ResponseSubmitFlightPlan](func() error {
		return flightPlanService.flightPlanOperation.SaveFlightPlan(req.FlightPlan)
	}); res != nil {
		return res
	}

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.FlushFlightPlan,
		Data: &fsd.FlushFlightPlan{
			TargetCallsign: req.FlightPlan.Callsign,
			TargetCid:      req.JwtHeader.Cid,
			FlightPlan:     req.FlightPlan,
		},
	})

	data := ResponseSubmitFlightPlan(true)
	return NewApiResponse(SuccessSubmitFlightPlan, &data)
}

func (flightPlanService *FlightPlanService) GetFlightPlan(req *RequestGetFlightPlan) *ApiResponse[ResponseGetFlightPlan] {
	flightPlan, res := CallDBFunc[*operation.FlightPlan, ResponseGetFlightPlan](func() (*operation.FlightPlan, error) {
		return flightPlanService.flightPlanOperation.GetFlightPlanByCid(req.Cid)
	})
	if res != nil {
		return res
	}

	return NewApiResponse(SuccessGetFlightPlan, &ResponseGetFlightPlan{FlightPlan: flightPlan})
}

func (flightPlanService *FlightPlanService) GetFlightPlans(req *RequestGetFlightPlans) *ApiResponse[ResponseGetFlightPlans] {
	if req.Page <= 0 || req.PageSize <= 0 {
		return NewApiResponse[ResponseGetFlightPlans](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseGetFlightPlans](req.Permission, operation.FlightPlanShowList); res != nil {
		return res
	}

	flightPlans, total, err := flightPlanService.flightPlanOperation.GetFlightPlans(req.Page, req.PageSize)
	if res := CheckDatabaseError[ResponseGetFlightPlans](err); res != nil {
		return res
	}

	return NewApiResponse(SuccessGetFlightPlans, &ResponseGetFlightPlans{
		Items:    flightPlans,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

func (flightPlanService *FlightPlanService) DeleteSelfFlightPlan(req *RequestDeleteSelfFlightPlan) *ApiResponse[ResponseDeleteSelfFlightPlan] {
	flightPlan, res := CallDBFunc[*operation.FlightPlan, ResponseDeleteSelfFlightPlan](func() (*operation.FlightPlan, error) {
		return flightPlanService.flightPlanOperation.GetFlightPlanByCid(req.Cid)
	})
	if res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseDeleteSelfFlightPlan](func() error {
		return flightPlanService.flightPlanOperation.DeleteSelfFlightPlan(flightPlan)
	}); res != nil {
		return res
	}

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: flightPlanService.auditLogOperation.NewAuditLog(
			operation.FlightPlanSelfDeleted,
			req.Cid,
			fmt.Sprintf("%04d", req.Cid),
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.FlushFlightPlan,
		Data: &fsd.FlushFlightPlan{
			TargetCallsign: flightPlan.Callsign,
			TargetCid:      req.Cid,
			FlightPlan:     nil,
		},
	})

	data := ResponseDeleteSelfFlightPlan(true)
	return NewApiResponse(SuccessDeleteSelfFlightPlan, &data)
}

func (flightPlanService *FlightPlanService) DeleteFlightPlan(req *RequestDeleteFlightPlan) *ApiResponse[ResponseDeleteFlightPlan] {
	if req.TargetCid <= 0 {
		return NewApiResponse[ResponseDeleteFlightPlan](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseDeleteFlightPlan](req.Permission, operation.FlightPlanDelete); res != nil {
		return res
	}

	flightPlan, res := CallDBFunc[*operation.FlightPlan, ResponseDeleteFlightPlan](func() (*operation.FlightPlan, error) {
		return flightPlanService.flightPlanOperation.GetFlightPlanByCid(req.TargetCid)
	})
	if res != nil {
		return res
	}

	if res := CallDBFuncWithoutRet[ResponseDeleteFlightPlan](func() error {
		return flightPlanService.flightPlanOperation.DeleteFlightPlan(flightPlan)
	}); res != nil {
		return res
	}

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: flightPlanService.auditLogOperation.NewAuditLog(
			operation.FlightPlanDeleted,
			req.Cid,
			fmt.Sprintf("%04d", req.TargetCid),
			req.Ip,
			req.UserAgent,
			nil,
		),
	})

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.FlushFlightPlan,
		Data: &fsd.FlushFlightPlan{
			TargetCallsign: flightPlan.Callsign,
			TargetCid:      req.Cid,
			FlightPlan:     nil,
		},
	})

	data := ResponseDeleteFlightPlan(true)
	return NewApiResponse(SuccessDeleteFlightPlan, &data)
}

func (flightPlanService *FlightPlanService) LockFlightPlan(req *RequestLockFlightPlan) *ApiResponse[ResponseLockFlightPlan] {
	if req.TargetCid <= 0 {
		return NewApiResponse[ResponseLockFlightPlan](ErrIllegalParam, nil)
	}

	if res := CheckPermission[ResponseLockFlightPlan](req.Permission, operation.FlightPlanChangeLock); res != nil {
		return res
	}

	flightPlan, res := CallDBFunc[*operation.FlightPlan, ResponseLockFlightPlan](func() (*operation.FlightPlan, error) {
		return flightPlanService.flightPlanOperation.GetFlightPlanByCid(req.TargetCid)
	})
	if res != nil {
		return res
	}

	if flightPlan.Locked == req.Lock {
		if req.Lock {
			return NewApiResponse[ResponseLockFlightPlan](ErrFlightPlanLocked, nil)
		} else {
			return NewApiResponse[ResponseLockFlightPlan](ErrFlightPlanUnlocked, nil)
		}
	}

	if res := CallDBFuncWithoutRet[ResponseLockFlightPlan](func() error {
		if req.Lock {
			return flightPlanService.flightPlanOperation.LockFlightPlan(flightPlan)
		} else {
			return flightPlanService.flightPlanOperation.UnlockFlightPlan(flightPlan)
		}
	}); res != nil {
		return res
	}

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.ChangeFlightPlanLockStatus,
		Data: &fsd.LockChange{
			TargetCallsign: flightPlan.Callsign,
			TargetCid:      req.Cid,
			Locked:         req.Lock,
		},
	})

	var auditLogType operation.AuditEventType
	if req.Lock {
		auditLogType = operation.FlightPlanLock
	} else {
		auditLogType = operation.FlightPlanUnlock
	}

	flightPlanService.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: flightPlanService.auditLogOperation.NewAuditLog(
			auditLogType,
			req.Cid,
			fmt.Sprintf("%04d", req.TargetCid),
			req.Ip,
			req.UserAgent,
			&operation.ChangeDetail{
				OldValue: strconv.FormatBool(!req.Lock),
				NewValue: strconv.FormatBool(req.Lock),
			},
		),
	})

	data := ResponseLockFlightPlan(true)
	return NewApiResponse(SuccessLockFlightPlan, &data)
}
