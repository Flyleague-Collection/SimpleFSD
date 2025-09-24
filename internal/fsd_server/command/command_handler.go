// Package command
package command

import (
	"errors"
	"fmt"
	c "github.com/half-nothing/simple-fsd/internal/fsd_server/client"
	"github.com/half-nothing/simple-fsd/internal/interfaces"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/global"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/interfaces/queue"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"strconv"
	"strings"
	"time"
)

// verifyUserInfo 验证用户信息与处理客户端重连机制
func (content *CommandContent) verifyUserInfo(session SessionInterface, callsign string, protocol int, cid operation.UserId, password string) *Result {
	if !callsignValid(callsign) {
		return ResultError(CallsignInvalid, true, callsign, nil)
	}

	if protocol != 9 {
		return ResultError(InvalidProtocolVision, true, callsign, nil)
	}

	client, ok := content.clientManager.GetClient(callsign)

	// 客户端存在且标记为断开连接
	if ok {
		if client.Reconnect(session) {
			// 客户端重连
			session.SetClient(client)
		} else {
			// 呼号已被使用
			return ResultError(CallsignInUse, true, callsign, nil)
		}
	}

	user, err := cid.GetUser(content.userOperation)
	if err != nil {
		return ResultError(InvalidCidPassword, true, callsign, err)
	}
	if user.Rating <= Ban.Index() {
		return ResultError(CidSuspended, true, callsign, nil)
	}
	if !content.userOperation.VerifyUserPassword(user, password) {
		return ResultError(InvalidCidPassword, true, callsign, nil)
	}

	// 重设重连客户端的User
	if client != nil {
		client.SetUser(user)
	}
	session.SetUser(user)

	return nil
}

func (content *CommandContent) checkRangeLimit(_ SessionInterface, realFacility Facility, realRange int) *Result {
	rangeLimit := realFacility.GetRangeLimit()
	if rangeLimit > -1 && realRange > rangeLimit {
		return ResultError(Custom, content.refuseOutRange, strconv.Itoa(realRange), fmt.Errorf("visual range out of limit, your visual range is %d but limit is %d", realRange, rangeLimit))
	}
	return nil
}

func (content *CommandContent) HandleAddAtc(session SessionInterface, data []string, rawLine []byte) *Result {
	// #AA 2352_OBS SERVER 2352 2352 123456  1  9  1  0  29.86379 119.49287 100
	// [0] [   1  ] [  2 ] [ 3] [ 4] [  5 ] [6][7][8][9] [  10  ] [   11  ] [12]
	callsign := data[0]
	cid := operation.GetUserId(data[3])
	password := data[4]
	protocol := utils.StrToInt(data[6], 0)
	result := content.verifyUserInfo(session, callsign, protocol, cid, password)
	if result != nil {
		return result
	}
	reqRating := utils.StrToInt(data[5], 0)
	if reqRating > session.User().Rating {
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
		session.SetFacilityIdent(facility)
	}
	realName := data[2]
	latitude := utils.StrToFloat(data[9], 0)
	longitude := utils.StrToFloat(data[10], 0)
	if session.Client() == nil {
		client := c.NewClient(content.application, callsign, Rating(reqRating), protocol, realName, session, true)
		_ = client.SetPosition(0, latitude, longitude)
		_ = content.clientManager.AddClient(client)
		session.SetClient(client)
	} else {
		session.Client().SetRating(Rating(reqRating))
		session.Client().SetRealName(realName)
	}
	session.Client().SendLine(MakePacket(ClientQuery, global.FSDServerName, callsign, "ATIS"))
	go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToClientInRange)
	session.Client().SendMotd()
	content.logger.InfoF("[%s] ATC login successfully", callsign)
	return ResultSuccess()
}

func (content *CommandContent) HandleAddPilot(session SessionInterface, data []string, rawLine []byte) *Result {
	//	#AP CES2352 SERVER 2352 123456  1   9  16  Half_nothing ZGHA
	//  [0] [  1  ] [  2 ] [ 3] [  4 ] [5] [6] [7] [       8       ]
	callsign := data[0]
	cid := operation.GetUserId(data[2])
	password := data[3]
	protocol := utils.StrToInt(data[5], 0)
	result := content.verifyUserInfo(session, callsign, protocol, cid, password)
	if result != nil {
		return result
	}
	reqRating := Rating(utils.StrToInt(data[4], 0) - 1)
	if reqRating != Normal || !RatingFacilityMap[reqRating].CheckFacility(Pilot) {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	simType := utils.StrToInt(data[6], 0)
	realName := data[7]
	if session.Client() == nil {
		client := c.NewClient(content.application, callsign, reqRating, protocol, realName, session, false)
		client.SetSimType(simType)
		session.SetClient(client)
		_ = content.clientManager.AddClient(client)
	} else {
		session.Client().SetRating(reqRating)
		session.Client().SetRealName(realName)
	}
	go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToClientInRange)
	session.Client().SendMotd()
	content.logger.InfoF("[%s] client login successfully", callsign)
	if !content.isSimulatorServer {
		flightPlan := session.Client().FlightPlan()
		if flightPlan != nil && flightPlan.FromWeb && callsign != flightPlan.Callsign {
			session.Client().SendLine(MakePacket(Message, "FPlanManager", callsign,
				fmt.Sprintf("Seems you are connect with callsign(%s), "+
					"but we found a flightplan submit by web at %s which has callsign(%s), "+
					"please check it.", callsign, flightPlan.UpdatedAt.Format(time.DateTime), flightPlan.Callsign)))
		}
	}
	return ResultSuccess()
}

func (content *CommandContent) HandleAtcPosUpdate(session SessionInterface, data []string, rawLine []byte) *Result {
	//  %  ZSHA_CTR 24550  6  600  5  27.28025 118.28701  0
	// [0] [   1  ] [ 2 ] [3] [4] [5] [   6  ] [   7   ] [8]
	callsign := data[0]
	rating := Rating(utils.StrToInt(data[4], 0))
	facility := Facility(1 << utils.StrToInt(data[2], 0))
	if !rating.CheckRatingFacility(facility) {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	if !session.FacilityIdent().CheckFacility(facility) {
		return ResultError(CallsignInvalid, true, callsign, errors.New("callsign and faility mismatch"))
	}
	if res := content.checkRangeLimit(session, facility, utils.StrToInt(data[3], 0)); res != nil {
		return res
	}
	frequency := utils.StrToInt(data[1], 0)
	visualRange := utils.StrToFloat(data[3], 0)
	latitude := utils.StrToFloat(data[5], 0)
	longitude := utils.StrToFloat(data[6], 0)
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToClientInRange)
	session.Client().UpdateAtcPos(frequency, facility, visualRange, latitude, longitude)
	return ResultSuccess()
}

func (content *CommandContent) HandlePilotPosUpdate(session SessionInterface, data []string, rawLine []byte) *Result {
	//	@   S  CPA421 7000  1  38.96244 121.53479 87   0  4290770974 278
	// [0] [1] [  2 ] [ 3] [4] [   5  ] [   6   ] [7] [8] [    9   ] [10]
	transponder := utils.StrToInt(data[2], 0)
	latitude := utils.StrToFloat(data[4], 0)
	longitude := utils.StrToFloat(data[5], 0)
	altitude := utils.StrToInt(data[6], 0)
	groundSpeed := utils.StrToInt(data[7], 0)
	pbh := uint32(utils.StrToInt(data[8], 0))
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToClientInRange)
	session.Client().UpdatePilotPos(transponder, latitude, longitude, altitude, groundSpeed, pbh)
	return ResultSuccess()
}

func (content *CommandContent) HandleAtcVisPointUpdate(session SessionInterface, data []string, _ []byte) *Result {
	//  '  ZSHA_CTR  0  36.67349 120.45621
	// [0] [   1  ] [2] [   3  ] [   4   ]
	visPos := utils.StrToInt(data[1], 0)
	latitude := utils.StrToFloat(data[2], 0)
	longitude := utils.StrToFloat(data[3], 0)
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	_ = session.Client().UpdateAtcVisPoint(visPos, latitude, longitude)
	return ResultSuccess()
}

// sendFrequencyMessage 发送频率消息
func (content *CommandContent) sendFrequencyMessage(session SessionInterface, targetStation string, rawLine []byte) *Result {
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	frequency := utils.StrToInt(targetStation[1:], -1)
	if frequency == -1 {
		return ResultError(Syntax, false, targetStation, fmt.Errorf("illegal frequency %s", targetStation))
	}
	if frequencyValid(frequency) {
		// 合法频率, 发给所有客户端
		go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToClientInRange)
	} else {
		// 非法频率, 大概率是管制使用, 只发给管制
		go content.clientManager.BroadcastMessage(rawLine, session.Client(), CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	}
	return nil
}

func (content *CommandContent) HandleClientQuery(session SessionInterface, data []string, rawLine []byte) *Result {
	//	查询飞行计划
	//	$CQ ZYSH_CTR SERVER FP  CPA421
	//  [0] [  1   ] [  2 ] [3] [  4 ]
	//
	//	修改飞行计划
	//	$CQ ZYSH_CTR @94835 FA  CPA421 31100
	//	[0] [  1   ] [  2 ] [3] [  4 ] [ 5 ]
	if session.Client() == nil {
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
			client, ok := content.clientManager.GetClient(targetCallsign)
			if !ok || client.FlightPlan() == nil {
				return ResultError(NoFlightPlan, false, session.Client().Callsign(), nil)
			}
			session.Client().SendLine([]byte(content.flightPlanOperation.ToString(client.FlightPlan(), data[0])))
		case AvailableAtc:
			available := isValidAtc(data[0])
			if available {
				session.Client().SendLine(MakePacket(ClientResponse, global.FSDServerName, data[0], "ATC:Y", data[0]))
			} else {
				session.Client().SendLine(MakePacket(ClientResponse, global.FSDServerName, data[0], "ATC:N", data[0]))
			}
		case ClientCapacity:
			if session.Client().IsAtc() {
				session.Client().SendLine(MakePacket(ClientResponse, global.FSDServerName, data[0], "CAPS:ATCINFO=1:SECPOS=1:FASTPOS=1:OBSPILOT=1"))
			} else {
				session.Client().SendLine(MakePacket(ClientResponse, global.FSDServerName, data[0], "CAPS:VERSION=1:ATCINFO=1:MODELDESC=1:ACCONFIG=1:VISUPDATE=1"))
			}
		case IpAddress:
			session.Client().SendLine(MakePacket(ClientResponse, global.FSDServerName, data[0], "IP", session.Conn().RemoteAddr().String()[0:strings.LastIndex(session.Conn().RemoteAddr().String(), ":")]))
		}
		return ResultSuccess()
	}
	// 如果发送目标是一个频率
	if strings.HasPrefix(targetStation, "@") {
		err := content.sendFrequencyMessage(session, targetStation, rawLine)
		// 如果目标频率是94835
		if targetStation == EuroscopeFrequency {
			subQuery := data[2]
			if !content.isSimulatorServer {
				if !session.Client().CheckFacility(AllowAtcFacility) {
					return ResultError(InvalidCtrl, false, session.Client().Callsign(), nil)
				}
				if subQuery == EditFlightPlan && commandLength >= 5 {
					targetCallsign := data[3]
					client, ok := content.clientManager.GetClient(targetCallsign)
					if !ok {
						// 这里并不是发给服务器的, 所以如果找不到指定客户端, 直接返回就行
						return ResultSuccess()
					}
					cruiseAltitude := utils.StrToInt(data[4], 0)
					if err := content.flightPlanOperation.UpdateCruiseAltitude(client.FlightPlan(), fmt.Sprintf("FL%03d", cruiseAltitude/100)); err != nil {
						// 这里并不是发给服务器的, 所以如果出错, 直接返回就行
						return ResultSuccess()
					}
				}
			}
			// ATIS信息更新, 需要更新服务器存储的ATIS信息
			if subQuery == InfoUpdate {
				session.Client().ClearAtcAtisInfo()
				session.Client().SendLine(MakePacket(ClientQuery, global.FSDServerName, session.Callsign(), AtcAtis))
			}
			if subQuery == Break {
				session.Client().SetBreak(true)
			}
			if subQuery == NoBreak {
				session.Client().SetBreak(false)
			}
		}
		if err != nil {
			return err
		}
	} else {
		_ = content.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

func (content *CommandContent) HandleClientResponse(session SessionInterface, data []string, rawLine []byte) *Result {
	//	$CR ZSHA_CTR ZSSS_APP CAPS ATCINFO=1 SECPOS=1 MODELDESC=1 ONGOINGCOORD=1 NEWINFO=1 TEAMSPEAK=1 ICAOEQ=1
	//  [0] [   1  ] [   2  ] [ 3] [   4   ] [  5   ] [    6    ] [     7      ] [   8   ] [     9   ] [  10  ]
	//	$CR ZSHA_CTR SERVER ATIS  T  ZSHA_CTR Shanghai Control
	//	[0] [   1  ] [  2 ] [ 3] [4] [           5           ]
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	commandLength := len(data)
	targetStation := data[1]
	if targetStation == global.FSDServerName {
		subQuery := data[2]
		if subQuery == "ATIS" && commandLength >= 5 {
			if data[3] == "T" {
				session.Client().AddAtcAtisInfo(data[4])
			}
			if data[3] == "Z" {
				session.Client().SetLogoffTime(data[4])
			}
		}
	}
	if strings.HasPrefix(targetStation, "@") {
		result := content.sendFrequencyMessage(session, targetStation, rawLine)
		if result != nil {
			return result
		}
	} else {
		_ = content.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

func (content *CommandContent) HandleMessage(session SessionInterface, data []string, rawLine []byte) *Result {
	// #TM ZSHA_CTR ZSSS_APP 111
	// [0] [   1  ] [   2  ] [3]
	targetStation := data[1]
	if strings.HasPrefix(targetStation, "@") {
		result := content.sendFrequencyMessage(session, targetStation, rawLine)
		if result != nil {
			return result
		}
		return ResultSuccess()
	}
	if strings.HasPrefix(targetStation, "*") {
		// 广播消息
		if targetStation == string(AllSup) {
			go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToSup)
			return ResultSuccess()
		}
		if targetStation == "*" && session.Client().IsAtc() && session.Client().CheckRating(AllowKillRating) {
			go content.clientManager.BroadcastMessage(rawLine, session.Client(), BroadcastToAll)
		}
		return ResultSuccess()
	}
	_ = content.clientManager.SendMessageTo(targetStation, rawLine)
	return ResultSuccess()
}

func (content *CommandContent) HandlePlan(session SessionInterface, data []string, rawLine []byte) *Result {
	// $FP CPA421 SERVER  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  1    18   2    26  ZYCC
	// [0] [  1 ] [  2 ] [3] [  4   ] [5] [ 6] [ 7] [8] [ 9 ] [10] [11] [12] [13] [14] [15]
	// /V/ SEL/AHFL VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [    16    ] [                      17                       ]
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if session.Client().IsAtc() {
		return ResultError(Syntax, false, "FLIGHT_PLAN", fmt.Errorf("atc can not submit fligth plan"))
	}
	if err := session.Client().UpsertFlightPlan(data); err != nil {
		return ResultError(Custom, false, "FLIGHT_PLAN", err)
	}
	if !session.Client().FlightPlan().Locked {
		go content.clientManager.BroadcastMessage(rawLine, session.Client(), CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	}
	return ResultSuccess()
}

func (content *CommandContent) HandleAtcEditPlan(session SessionInterface, data []string, _ []byte) *Result {
	// $AM ZYSH_CTR SERVER CPA421  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  11  8     22   6   ZYCC
	// [0] [   1  ] [  2 ] [  3 ] [4] [   5  ] [6] [ 7] [ 8] [9] [ 10] [11] [12] [13] [14] [15] [16]
	// /V/ SEL/AHFL CHI19D/28 VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [     17   ] [                             18                          ]
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if !session.Client().IsAtc() {
		return ResultError(Syntax, false, session.Client().Callsign(), fmt.Errorf("only act can edit flight plan"))
	}
	if !session.Client().CheckFacility(AllowAtcFacility) {
		return ResultError(Syntax, false, session.Client().Callsign(), fmt.Errorf("%s facility not allowed to edit plan", session.Client().Facility().String()))
	}
	targetCallsign := data[2]
	client, ok := content.clientManager.GetClient(targetCallsign)
	if !ok {
		return ResultError(NoCallsignFound, false, session.Client().Callsign(), fmt.Errorf("%s not exists", targetCallsign))
	}
	if client.FlightPlan == nil {
		return ResultError(NoFlightPlan, false, session.Client().Callsign(), fmt.Errorf("%s do not have filght plan", session.Client().Callsign()))
	}
	client.FlightPlan().Locked = !content.isSimulatorServer
	if err := content.flightPlanOperation.UpdateFlightPlan(client.FlightPlan(), data[1:], true); err != nil {
		return ResultError(Syntax, false, session.Client().Callsign(), err)
	}
	go content.clientManager.BroadcastMessage([]byte(content.flightPlanOperation.ToString(client.FlightPlan(), string(AllATC))),
		session.Client(), CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	return ResultSuccess()
}

func (content *CommandContent) HandleKillClient(session SessionInterface, data []string, _ []byte) *Result {
	// $!! ZSHA_CTR CPA421 test
	if session.Client() == nil {
		return ResultError(Custom, false, "", fmt.Errorf("client not register"))
	}
	if !(session.Client().IsAtc() && session.Client().CheckRating(AllowKillRating)) {
		return ResultError(Custom, false, session.Client().Callsign(), fmt.Errorf("%s rating not allowed to kill client", session.Client().Rating().String()))
	}
	targetStation := data[1]
	client, err := content.clientManager.KickClientFromServer(targetStation, data[2])
	if err != nil {
		return ResultError(NoCallsignFound, false, session.Client().Callsign(), fmt.Errorf("%s not exists", targetStation))
	}
	content.messageQueue.Publish(&queue.Message{
		Type: queue.SendKickedFromServerEmail,
		Data: &interfaces.KickedFromServerEmailData{
			User:     client.User(),
			Operator: session.User(),
			Reason:   data[2],
		},
	})
	content.messageQueue.Publish(&queue.Message{
		Type: queue.AuditLog,
		Data: content.auditLogOperation.NewAuditLog(
			operation.ClientKickedFsd,
			session.User().Cid,
			fmt.Sprintf("%s(%04d)", client.Callsign(), client.User().Cid),
			session.ConnId(),
			"NOT AVAILABLE",
			nil,
		),
	})
	return ResultSuccess()
}

func (content *CommandContent) HandleRequest(_ SessionInterface, data []string, rawLine []byte) *Result {
	targetStation := data[1]
	_ = content.clientManager.SendMessageTo(targetStation, rawLine)
	return ResultSuccess()
}

func (content *CommandContent) RemoveClient(session SessionInterface, _ []string, _ []byte) *Result {
	content.logger.InfoF("[%s] Offline", session.Client().Callsign())
	return ResultSuccess()
}

func (content *CommandContent) HandleSquawkBox(session SessionInterface, data []string, rawLine []byte) *Result {
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	targetStation := data[1]
	client, ok := content.clientManager.GetClient(targetStation)
	if !ok {
		return ResultError(NoCallsignFound, false, session.Client().Callsign(), fmt.Errorf("%s not exists", targetStation))
	}
	client.SendLine(rawLine)
	return ResultSuccess()
}

func (content *CommandContent) HandleWeatherQuery(session SessionInterface, data []string, _ []byte) *Result {
	if session.Client() == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	targetStation := data[3]
	if len(targetStation) != 4 {
		return ResultError(Syntax, false, targetStation, fmt.Errorf("invalid target station"))
	}
	result, err := content.metarManager.QueryMetar(targetStation)
	if err != nil {
		return ResultError(NoWeatherProfile, false, targetStation, fmt.Errorf("cant fetch metar for %s", targetStation))
	} else {
		session.Client().SendLine(MakePacket(WeatherResponse, global.FSDServerName, session.Callsign(), "METAR", result[6:]))
	}
	return ResultSuccess()
}
