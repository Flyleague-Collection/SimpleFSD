// Package packet 命令处理的核心函数定义
package packet

import (
	"errors"
	"fmt"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	"github.com/half-nothing/simple-fsd/internal/interfaces/service"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"strconv"
	"strings"
	"time"
)

func (session *Session) checkPacketLength(data []string, requirement *CommandRequirement) (*Result, bool) {
	length := len(data)
	if length < requirement.RequireLength {
		return ResultError(Syntax, requirement.Fatal, session.callsign, fmt.Errorf("datapack length too short, require %d but got %d", requirement.RequireLength, length)), false
	}
	return nil, true
}

// verifyUserInfo 验证用户信息与处理客户端重连机制
func (session *Session) verifyUserInfo(callsign string, protocol int, cid UserId, password string) *Result {
	if !callsignValid(callsign) {
		return ResultError(CallsignInvalid, true, callsign, nil)
	}

	if protocol != 9 {
		return ResultError(InvalidProtocolVision, true, callsign, nil)
	}

	client, ok := session.clientManager.GetClient(callsign)

	// 客户端存在且标记为断开连接
	if ok {
		if client.Reconnect(session) {
			// 客户端重连
			session.client = client
		} else {
			// 呼号已被使用
			return ResultError(CallsignInUse, true, callsign, nil)
		}
	}

	user, err := cid.GetUser(session.userOperation)
	if err != nil {
		return ResultError(InvalidCidPassword, true, callsign, err)
	}
	if user.Rating <= Ban.Index() {
		return ResultError(CidSuspended, true, callsign, nil)
	}
	if !session.userOperation.VerifyUserPassword(user, password) {
		return ResultError(InvalidCidPassword, true, callsign, nil)
	}

	// 重设重连客户端的User
	if client != nil {
		client.SetUser(user)
	}
	session.user = user

	return nil
}

func (session *Session) checkRangeLimit(realFacility Facility, realRange int) *Result {
	rangeLimit := realFacility.GetRangeLimit()
	if rangeLimit > -1 && realRange > rangeLimit {
		return ResultError(Custom, session.refuseOutRange, strconv.Itoa(realRange), fmt.Errorf("visual range out of limit, your visual range is %d but limit is %d", realRange, rangeLimit))
	}
	return nil
}

// handleAddAtc 处理管制员登录
func (session *Session) handleAddAtc(data []string, rawLine []byte) *Result {
	// #AA 2352_OBS SERVER 2352 2352 123456  1  9  1  0  29.86379 119.49287 100
	// [0] [   1  ] [  2 ] [ 3] [ 4] [  5 ] [6][7][8][9] [  10  ] [   11  ] [12]
	callsign := data[0]
	cid := GetUserId(data[3])
	password := data[4]
	protocol := utils.StrToInt(data[6], 0)
	result := session.verifyUserInfo(callsign, protocol, cid, password)
	if result != nil {
		return result
	}
	reqRating := utils.StrToInt(data[5], 0)
	if reqRating > session.user.Rating {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	facilityIdent := strings.Split(callsign, "_")
	if len(facilityIdent) < 2 {
		return ResultError(Custom, true, callsign, fmt.Errorf("invalid callsign %s", callsign))
	}
	ident := facilityIdent[len(facilityIdent)-1]
	if facility, exist := FacilityMap[ident]; !exist {
		return ResultError(Custom, true, callsign, fmt.Errorf("invalid callsign %s", callsign))
	} else {
		session.facilityIdent = facility
	}
	realName := data[2]
	latitude := utils.StrToFloat(data[9], 0)
	longitude := utils.StrToFloat(data[10], 0)
	if session.client == nil {
		session.client = NewClient(session.application, callsign, Rating(reqRating), protocol, realName, session, true)
		_ = session.client.SetPosition(0, latitude, longitude)
		_ = session.clientManager.AddClient(session.client)
	} else {
		session.client.SetRating(Rating(reqRating))
		session.client.SetRealName(realName)
	}
	session.client.SendLine(makePacket(ClientQuery, global.FSDServerName, callsign, "ATIS"))
	go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToClientInRange)
	session.client.SendMotd()
	session.logger.InfoF("[%s] ATC login successfully", callsign)
	return ResultSuccess()
}

// handleAddPilot 处理客户端登录
func (session *Session) handleAddPilot(data []string, rawLine []byte) *Result {
	//	#AP CES2352 SERVER 2352 123456  1   9  16  Half_nothing ZGHA
	//  [0] [  1  ] [  2 ] [ 3] [  4 ] [5] [6] [7] [       8       ]
	callsign := data[0]
	cid := GetUserId(data[2])
	password := data[3]
	protocol := utils.StrToInt(data[5], 0)
	result := session.verifyUserInfo(callsign, protocol, cid, password)
	if result != nil {
		return result
	}
	reqRating := Rating(utils.StrToInt(data[4], 0) - 1)
	if reqRating != Normal || !RatingFacilityMap[reqRating].CheckFacility(Pilot) {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	simType := utils.StrToInt(data[6], 0)
	realName := data[7]
	if session.client == nil {
		session.client = NewClient(session.application, callsign, reqRating, protocol, realName, session, false)
		session.client.SetSimType(simType)
		_ = session.clientManager.AddClient(session.client)
	} else {
		session.client.SetRating(reqRating)
		session.client.SetRealName(realName)
	}
	go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToClientInRange)
	session.client.SendMotd()
	session.logger.InfoF("[%s] client login successfully", callsign)
	if !session.isSimulatorServer {
		flightPlan := session.client.FlightPlan()
		if flightPlan != nil && flightPlan.FromWeb && callsign != flightPlan.Callsign {
			session.client.SendLine(makePacket(Message, "FPlanManager", callsign,
				fmt.Sprintf("Seems you are connect with callsign(%s), "+
					"but we found a flightplan submit by web at %s which has callsign(%s), "+
					"please check it.", callsign, flightPlan.UpdatedAt.Format(time.DateTime), flightPlan.Callsign)))
		}
	}
	return ResultSuccess()
}

// handleAtcPosUpdate 处理管制员位置更新
func (session *Session) handleAtcPosUpdate(data []string, rawLine []byte) *Result {
	//  %  ZSHA_CTR 24550  6  600  5  27.28025 118.28701  0
	// [0] [   1  ] [ 2 ] [3] [4] [5] [   6  ] [   7   ] [8]
	callsign := data[0]
	rating := Rating(utils.StrToInt(data[4], 0))
	facility := Facility(1 << utils.StrToInt(data[2], 0))
	if !rating.CheckRatingFacility(facility) {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	if !session.facilityIdent.CheckFacility(facility) {
		return ResultError(CallsignInvalid, true, callsign, errors.New("callsign and faility mismatch"))
	}
	if res := session.checkRangeLimit(facility, utils.StrToInt(data[3], 0)); res != nil {
		return res
	}
	frequency := utils.StrToInt(data[1], 0)
	visualRange := utils.StrToFloat(data[3], 0)
	latitude := utils.StrToFloat(data[5], 0)
	longitude := utils.StrToFloat(data[6], 0)
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToClientInRange)
	session.client.UpdateAtcPos(frequency, facility, visualRange, latitude, longitude)
	return ResultSuccess()
}

// handlePilotPosUpdate 处理飞行员位置更新
func (session *Session) handlePilotPosUpdate(data []string, rawLine []byte) *Result {
	//	@   S  CPA421 7000  1  38.96244 121.53479 87   0  4290770974 278
	// [0] [1] [  2 ] [ 3] [4] [   5  ] [   6   ] [7] [8] [    9   ] [10]
	transponder := utils.StrToInt(data[2], 0)
	latitude := utils.StrToFloat(data[4], 0)
	longitude := utils.StrToFloat(data[5], 0)
	altitude := utils.StrToInt(data[6], 0)
	groundSpeed := utils.StrToInt(data[7], 0)
	pbh := uint32(utils.StrToInt(data[8], 0))
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToClientInRange)
	session.client.UpdatePilotPos(transponder, latitude, longitude, altitude, groundSpeed, pbh)
	return ResultSuccess()
}

// handleAtcVisPointUpdate 处理管制员视程点更新
func (session *Session) handleAtcVisPointUpdate(data []string, _ []byte) *Result {
	//  '  ZSHA_CTR  0  36.67349 120.45621
	// [0] [   1  ] [2] [   3  ] [   4   ]
	visPos := utils.StrToInt(data[1], 0)
	latitude := utils.StrToFloat(data[2], 0)
	longitude := utils.StrToFloat(data[3], 0)
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	_ = session.client.UpdateAtcVisPoint(visPos, latitude, longitude)
	return ResultSuccess()
}

// sendFrequencyMessage 发送频率消息
func (session *Session) sendFrequencyMessage(targetStation string, rawLine []byte) *Result {
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	frequency := utils.StrToInt(targetStation[1:], -1)
	if frequency == -1 {
		return ResultError(Syntax, false, targetStation, fmt.Errorf("illegal frequency %s", targetStation))
	}
	if frequencyValid(frequency) {
		// 合法频率, 发给所有客户端
		go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToClientInRange)
	} else {
		// 非法频率, 大概率是管制使用, 只发给管制
		go session.clientManager.BroadcastMessage(rawLine, session.client, CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	}
	return nil
}

// handleClientQuery 处理客户端查询消息
func (session *Session) handleClientQuery(data []string, rawLine []byte) *Result {
	//	查询飞行计划
	//	$CQ ZYSH_CTR SERVER FP  CPA421
	//  [0] [  1   ] [  2 ] [3] [  4 ]
	//
	//	修改飞行计划
	//	$CQ ZYSH_CTR @94835 FA  CPA421 31100
	//	[0] [  1   ] [  2 ] [3] [  4 ] [ 5 ]
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	commandLength := len(data)
	if commandLength < 3 {
		return ResultError(Syntax, false, "", fmt.Errorf("illegal command length %d", commandLength))
	}
	targetStation := data[1]
	if targetStation == global.FSDServerName {
		subQuery := data[2]
		// 查询指定机组的飞行计划
		switch subQuery {
		case ClientFlightPlan:
			if commandLength < 4 {
				return ResultError(Syntax, false, "", fmt.Errorf("illegal command length %d", commandLength))
			}
			targetCallsign := data[3]
			client, ok := session.clientManager.GetClient(targetCallsign)
			if !ok || client.FlightPlan() == nil {
				return ResultError(NoFlightPlan, false, session.client.Callsign(), nil)
			}
			session.client.SendLine([]byte(session.flightPlanOperation.ToString(client.FlightPlan(), data[0])))
		case AvailableAtc:
			available := isValidAtc(data[0])
			if available {
				session.client.SendLine(makePacket(ClientResponse, global.FSDServerName, data[0], "ATC:Y", data[0]))
			} else {
				session.client.SendLine(makePacket(ClientResponse, global.FSDServerName, data[0], "ATC:N", data[0]))
			}
		case ClientCapacity:
			if session.client.IsAtc() {
				session.client.SendLine(makePacket(ClientResponse, global.FSDServerName, data[0], "CAPS:ATCINFO=1:SECPOS=1:FASTPOS=1:OBSPILOT=1"))
			} else {
				session.client.SendLine(makePacket(ClientResponse, global.FSDServerName, data[0], "CAPS:VERSION=1:ATCINFO=1:MODELDESC=1:ACCONFIG=1:VISUPDATE=1"))
			}
		case IpAddress:
			session.client.SendLine(makePacket(ClientResponse, global.FSDServerName, data[0], "IP", session.conn.RemoteAddr().String()[0:strings.LastIndex(session.conn.RemoteAddr().String(), ":")]))
		}
		return ResultSuccess()
	}
	// 如果发送目标是一个频率
	if strings.HasPrefix(targetStation, "@") {
		err := session.sendFrequencyMessage(targetStation, rawLine)
		// 如果目标频率是94835
		if targetStation == EuroscopeFrequency {
			subQuery := data[2]
			if !session.isSimulatorServer {
				if !session.client.CheckFacility(AllowAtcFacility) {
					return ResultError(InvalidCtrl, false, session.client.Callsign(), nil)
				}
				if subQuery == EditFlightPlan && commandLength >= 5 {
					targetCallsign := data[3]
					client, ok := session.clientManager.GetClient(targetCallsign)
					if !ok {
						// 这里并不是发给服务器的, 所以如果找不到指定客户端, 直接返回就行
						return ResultSuccess()
					}
					cruiseAltitude := utils.StrToInt(data[4], 0)
					if err := session.flightPlanOperation.UpdateCruiseAltitude(client.FlightPlan(), fmt.Sprintf("FL%03d", cruiseAltitude/100)); err != nil {
						// 这里并不是发给服务器的, 所以如果出错, 直接返回就行
						return ResultSuccess()
					}
				}
			}
			// ATIS信息更新, 需要更新服务器存储的ATIS信息
			if subQuery == InfoUpdate {
				session.client.ClearAtcAtisInfo()
				session.client.SendLine(makePacket(ClientQuery, global.FSDServerName, session.callsign, AtcAtis))
			}
			if subQuery == Break {
				session.client.SetBreak(true)
			}
			if subQuery == NoBreak {
				session.client.SetBreak(false)
			}
		}
		if err != nil {
			return err
		}
	} else {
		_ = session.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

// handleClientResponse 处理客户端回复消息
func (session *Session) handleClientResponse(data []string, rawLine []byte) *Result {
	//	$CR ZSHA_CTR ZSSS_APP CAPS ATCINFO=1 SECPOS=1 MODELDESC=1 ONGOINGCOORD=1 NEWINFO=1 TEAMSPEAK=1 ICAOEQ=1
	//  [0] [   1  ] [   2  ] [ 3] [   4   ] [  5   ] [    6    ] [     7      ] [   8   ] [     9   ] [  10  ]
	//	$CR ZSHA_CTR SERVER ATIS  T  ZSHA_CTR Shanghai Control
	//	[0] [   1  ] [  2 ] [ 3] [4] [           5           ]
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	commandLength := len(data)
	targetStation := data[1]
	if targetStation == global.FSDServerName {
		subQuery := data[2]
		if subQuery == "ATIS" && commandLength >= 5 {
			if data[3] == "T" {
				session.client.AddAtcAtisInfo(data[4])
			}
			if data[3] == "Z" {
				session.client.SetLogoffTime(data[4])
			}
		}
	}
	if strings.HasPrefix(targetStation, "@") {
		result := session.sendFrequencyMessage(targetStation, rawLine)
		if result != nil {
			return result
		}
	} else {
		_ = session.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

func (session *Session) handleMessage(data []string, rawLine []byte) *Result {
	// #TM ZSHA_CTR ZSSS_APP 111
	// [0] [   1  ] [   2  ] [3]
	targetStation := data[1]
	if strings.HasPrefix(targetStation, "@") {
		result := session.sendFrequencyMessage(targetStation, rawLine)
		if result != nil {
			return result
		}
		return ResultSuccess()
	}
	if strings.HasPrefix(targetStation, "*") {
		// 广播消息
		if targetStation == string(AllSup) {
			go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToSup)
			return ResultSuccess()
		}
		if targetStation == "*" && session.client.IsAtc() && session.client.CheckRating(AllowKillRating) {
			go session.clientManager.BroadcastMessage(rawLine, session.client, BroadcastToAll)
		}
		return ResultSuccess()
	}
	_ = session.clientManager.SendMessageTo(targetStation, rawLine)
	return ResultSuccess()
}

func (session *Session) handlePlan(data []string, rawLine []byte) *Result {
	// $FP CPA421 SERVER  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  1    18   2    26  ZYCC
	// [0] [  1 ] [  2 ] [3] [  4   ] [5] [ 6] [ 7] [8] [ 9 ] [10] [11] [12] [13] [14] [15]
	// /V/ SEL/AHFL VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [    16    ] [                      17                       ]
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if session.client.IsAtc() {
		return ResultError(Syntax, false, "FLIGHT_PLAN", fmt.Errorf("atc can not submit fligth plan"))
	}
	if err := session.client.UpsertFlightPlan(data); err != nil {
		return ResultError(Custom, false, "FLIGHT_PLAN", err)
	}
	if !session.client.FlightPlan().Locked {
		go session.clientManager.BroadcastMessage(rawLine, session.client, CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	}
	return ResultSuccess()
}

func (session *Session) handleAtcEditPlan(data []string, _ []byte) *Result {
	// $AM ZYSH_CTR SERVER CPA421  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  11  8     22   6   ZYCC
	// [0] [   1  ] [  2 ] [  3 ] [4] [   5  ] [6] [ 7] [ 8] [9] [ 10] [11] [12] [13] [14] [15] [16]
	// /V/ SEL/AHFL CHI19D/28 VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [     17   ] [                             18                          ]
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if !session.client.IsAtc() {
		return ResultError(Syntax, false, session.client.Callsign(), fmt.Errorf("only act can edit flight plan"))
	}
	if !session.client.CheckFacility(AllowAtcFacility) {
		return ResultError(Syntax, false, session.client.Callsign(), fmt.Errorf("%s facility not allowed to edit plan", session.client.Facility().String()))
	}
	targetCallsign := data[2]
	client, ok := session.clientManager.GetClient(targetCallsign)
	if !ok {
		return ResultError(NoCallsignFound, false, session.client.Callsign(), fmt.Errorf("%s not exists", targetCallsign))
	}
	if client.FlightPlan == nil {
		return ResultError(NoFlightPlan, false, session.client.Callsign(), fmt.Errorf("%s do not have filght plan", session.client.Callsign()))
	}
	client.FlightPlan().Locked = !session.isSimulatorServer
	if err := session.flightPlanOperation.UpdateFlightPlan(client.FlightPlan(), data[1:], true); err != nil {
		return ResultError(Syntax, false, session.client.Callsign(), err)
	}
	go session.clientManager.BroadcastMessage([]byte(session.flightPlanOperation.ToString(client.FlightPlan(), string(AllATC))),
		session.client, CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	return ResultSuccess()
}

func (session *Session) handleKillClient(data []string, _ []byte) *Result {
	// $!! ZSHA_CTR CPA421 test
	if session.client == nil {
		return ResultError(Custom, false, "", fmt.Errorf("client not register"))
	}
	if !(session.client.IsAtc() && session.client.CheckRating(AllowKillRating)) {
		return ResultError(Custom, false, session.client.Callsign(), fmt.Errorf("%s rating not allowed to kill client", session.client.Rating().String()))
	}
	targetStation := data[1]
	client, err := session.clientManager.KickClientFromServer(targetStation, data[2])
	if err != nil {
		return ResultError(NoCallsignFound, false, session.client.Callsign(), fmt.Errorf("%s not exists", targetStation))
	}
	session.application.MessageQueue().Publish(&queue.Message{
		Type: queue.SendKickedFromServerEmail,
		Data: &service.KickedFromServerEmailData{
			User:     client.User(),
			Operator: session.user,
			Reason:   data[2],
		},
	})
	session.application.MessageQueue().Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: session.application.Operations().AuditLogOperation().NewAuditLog(
			ClientKickedFsd, session.user.Cid, fmt.Sprintf("%s(%04d)", client.Callsign(), client.User().Cid),
			session.connId, "NOT AVAILABLE", nil,
		),
	})
	return ResultSuccess()
}

func (session *Session) handleRequest(data []string, rawLine []byte) *Result {
	// $HO ZSSS_APP ZSHA_CTR CES2352
	// [0] [  1   ] [   2  ] [  3  ]
	// #PC ZSHA_CTR ZSSS_APP CCP HC CES2352
	// $HA ZSHA_CTR ZSSS_APP CES2352
	targetStation := data[1]
	_ = session.clientManager.SendMessageTo(targetStation, rawLine)
	return ResultSuccess()
}

func (session *Session) removeClient(_ []string, _ []byte) *Result {
	// #DA ZGGG_CTR SERVER
	session.logger.InfoF("[%s] Offline", session.client.Callsign())
	return ResultSuccess()
}

func (session *Session) handleSquawkBox(data []string, rawLine []byte) *Result {
	//# SB ZSHA_CTR CES7199 FSIPIR 0 ZZZ C172 3.73453 1.32984 174.00000 3.3028DC0.98CF9A41 L1P Cessna Skyhawk 172SP
	if session.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	targetStation := data[1]
	client, ok := session.clientManager.GetClient(targetStation)
	if !ok {
		return ResultError(NoCallsignFound, false, session.client.Callsign(), fmt.Errorf("%s not exists", targetStation))
	}
	client.SendLine(rawLine)
	return ResultSuccess()
}

func (session *Session) handleCommand(commandType ClientCommand, data []string, rawLine []byte) *Result {
	var result = ResultSuccess()
	if requirement, ok := CommandRequirements[commandType]; ok {
		if err, ok := session.checkPacketLength(data, requirement); !ok {
			return err
		}
	}
	switch commandType {
	case AddAtc:
		result = session.handleAddAtc(data, rawLine)
	case AddPilot:
		result = session.handleAddPilot(data, rawLine)
	case RequestHandoff, AcceptHandoff, ProController:
		result = session.handleRequest(data, rawLine)
	case PilotPosition:
		result = session.handlePilotPosUpdate(data, rawLine)
	case Plan:
		result = session.handlePlan(data, rawLine)
	case AtcEditPlan:
		result = session.handleAtcEditPlan(data, rawLine)
	case KillClient:
		result = session.handleKillClient(data, rawLine)
	case AtcPosition:
		result = session.handleAtcPosUpdate(data, rawLine)
	case AtcSubVisPoint:
		result = session.handleAtcVisPointUpdate(data, rawLine)
	case Message:
		result = session.handleMessage(data, rawLine)
	case ClientQuery:
		result = session.handleClientQuery(data, rawLine)
	case ClientResponse:
		result = session.handleClientResponse(data, rawLine)
	case RemoveAtc, RemovePilot:
		session.removeClient(data, rawLine)
	case SquawkBox:
		result = session.handleSquawkBox(data, rawLine)
	default:
		result = ResultSuccess()
	}
	return result
}
