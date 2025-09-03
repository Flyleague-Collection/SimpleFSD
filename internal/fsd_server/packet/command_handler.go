// Package packet 命令处理的核心函数定义
package packet

import (
	"fmt"
	c "github.com/half-nothing/simple-fsd/internal/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"strings"
)

func (ch *ConnectionHandler) checkPacketLength(data []string, requirement *CommandRequirement) (*Result, bool) {
	length := len(data)
	if length < requirement.RequireLength {
		return ResultError(Syntax, requirement.Fatal, ch.callsign, fmt.Errorf("datapack length too short, require %d but got %d", requirement.RequireLength, length)), false
	}
	return nil, true
}

// verifyUserInfo 验证用户信息与处理客户端重连机制
func (ch *ConnectionHandler) verifyUserInfo(callsign string, protocol int, cid string, password string) *Result {
	if !callsignValid(callsign) {
		return ResultError(CallsignInvalid, true, callsign, nil)
	}

	if protocol != 9 {
		return ResultError(InvalidProtocolVision, true, callsign, nil)
	}

	client, ok := ch.clientManager.GetClient(callsign)

	// 客户端存在且标记为断开连接
	if ok {
		if client.Reconnect(ch) {
			// 客户端重连
			ch.client = client
		} else {
			// 呼号已被使用
			return ResultError(CallsignInUse, true, callsign, nil)
		}
	}

	user, err := ch.userOperation.GetUserByCid(cid)
	if err != nil {
		return ResultError(AuthFail, true, callsign, err)
	}
	if user.Rating == Ban.Index() {
		return ResultError(UserBaned, true, callsign, nil)
	}
	if !ch.userOperation.VerifyUserPassword(user, password) {
		return ResultError(AuthFail, true, callsign, nil)
	}

	// 重设重连客户端的User
	if client != nil {
		client.SetUser(user)
	}
	ch.user = user

	return nil
}

// handleAddAtc 处理管制员登录
func (ch *ConnectionHandler) handleAddAtc(data []string, rawLine []byte) *Result {
	// #AA 2352_OBS SERVER 2352 2352 123456  1  9  1  0  29.86379 119.49287 100
	// [0] [   1  ] [  2 ] [ 3] [ 4] [  5 ] [6][7][8][9] [  10  ] [   11  ] [12]
	callsign := data[0]
	cid := data[3]
	password := data[4]
	protocol := utils.StrToInt(data[6], 0)
	result := ch.verifyUserInfo(callsign, protocol, cid, password)
	if result != nil {
		return result
	}
	reqRating := utils.StrToInt(data[5], 0)
	if reqRating > ch.user.Rating {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	realName := data[2]
	latitude := utils.StrToFloat(data[9], 0)
	longitude := utils.StrToFloat(data[10], 0)
	if ch.client == nil {
		ch.client = ch.clientManager.NewClient(callsign, Rating(reqRating), protocol, realName, ch, true)
		_ = ch.client.SetPosition(0, latitude, longitude)
		_ = ch.clientManager.AddClient(ch.client)
	}
	ch.client.SendLine(makePacket(ClientQuery, "SERVER", callsign, "ATIS"))
	go ch.clientManager.BroadcastMessage(rawLine, ch.client, BroadcastToClientInRange)
	ch.client.SendMotd()
	c.InfoF("[%s] ATC login successfully", callsign)
	return ResultSuccess()
}

// handleAddPilot 处理客户端登录
func (ch *ConnectionHandler) handleAddPilot(data []string, rawLine []byte) *Result {
	//	#AP CES2352 SERVER 2352 123456  1   9  16  Half_nothing ZGHA
	//  [0] [  1  ] [  2 ] [ 3] [  4 ] [5] [6] [7] [       8       ]
	callsign := data[0]
	cid := data[2]
	password := data[3]
	protocol := utils.StrToInt(data[5], 0)
	result := ch.verifyUserInfo(callsign, protocol, cid, password)
	if result != nil {
		return result
	}
	reqRating := Rating(utils.StrToInt(data[4], 0) - 1)
	if reqRating != Normal || !RatingFacilityMap[reqRating].CheckFacility(Pilot) {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	simType := utils.StrToInt(data[6], 0)
	realName := data[7]
	if ch.client == nil {
		ch.client = ch.clientManager.NewClient(callsign, reqRating, protocol, realName, ch, false)
		ch.client.SetSimType(simType)
		_ = ch.clientManager.AddClient(ch.client)
	}
	go ch.clientManager.BroadcastMessage(rawLine, ch.client, BroadcastToClientInRange)
	ch.client.SendMotd()
	c.InfoF("[%s] client login successfully", callsign)
	return ResultSuccess()
}

// handleAtcPosUpdate 处理管制员位置更新
func (ch *ConnectionHandler) handleAtcPosUpdate(data []string, rawLine []byte) *Result {
	//  %  ZSHA_CTR 24550  6  600  5  27.28025 118.28701  0
	// [0] [   1  ] [ 2 ] [3] [4] [5] [   6  ] [   7   ] [8]
	callsign := data[0]
	facility := Facility(1 << utils.StrToInt(data[2], 0))
	rating := Rating(utils.StrToInt(data[4], 0))
	if !rating.CheckRatingFacility(facility) {
		return ResultError(RequestLevelTooHigh, true, callsign, nil)
	}
	frequency := utils.StrToInt(data[1], 0)
	visualRange := utils.StrToFloat(data[3], 0)
	latitude := utils.StrToFloat(data[5], 0)
	longitude := utils.StrToFloat(data[6], 0)
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	go ch.clientManager.BroadcastMessage(rawLine, ch.client, BroadcastToClientInRange)
	ch.client.UpdateAtcPos(frequency, facility, visualRange, latitude, longitude)
	return ResultSuccess()
}

// handlePilotPosUpdate 处理飞行员位置更新
func (ch *ConnectionHandler) handlePilotPosUpdate(data []string, rawLine []byte) *Result {
	//	@   S  CPA421 7000  1  38.96244 121.53479 87   0  4290770974 278
	// [0] [1] [  2 ] [ 3] [4] [   5  ] [   6   ] [7] [8] [    9   ] [10]
	transponder := utils.StrToInt(data[2], 0)
	latitude := utils.StrToFloat(data[4], 0)
	longitude := utils.StrToFloat(data[5], 0)
	altitude := utils.StrToInt(data[6], 0)
	groundSpeed := utils.StrToInt(data[7], 0)
	pbh := uint32(utils.StrToInt(data[8], 0))
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	go ch.clientManager.BroadcastMessage(rawLine, ch.client, BroadcastToClientInRange)
	ch.client.UpdatePilotPos(transponder, latitude, longitude, altitude, groundSpeed, pbh)
	return ResultSuccess()
}

// handleAtcVisPointUpdate 处理管制员视程点更新
func (ch *ConnectionHandler) handleAtcVisPointUpdate(data []string, _ []byte) *Result {
	//  '  ZSHA_CTR  0  36.67349 120.45621
	// [0] [   1  ] [2] [   3  ] [   4   ]
	visPos := utils.StrToInt(data[1], 0)
	latitude := utils.StrToFloat(data[2], 0)
	longitude := utils.StrToFloat(data[3], 0)
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	_ = ch.client.UpdateAtcVisPoint(visPos, latitude, longitude)
	return ResultSuccess()
}

// sendFrequencyMessage 发送频率消息
func (ch *ConnectionHandler) sendFrequencyMessage(targetStation string, rawLine []byte) *Result {
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	frequency := utils.StrToInt(targetStation[1:], -1)
	if frequency == -1 {
		return ResultError(Syntax, false, ch.client.Callsign(), fmt.Errorf("illegal frequency %s", targetStation))
	}
	if frequencyValid(frequency) {
		// 合法频率, 发给所有客户端
		go ch.clientManager.BroadcastMessage(rawLine, ch.client, BroadcastToClientInRange)
	} else {
		// 非法频率, 大概率是管制使用, 只发给管制
		go ch.clientManager.BroadcastMessage(rawLine, ch.client, CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	}
	return nil
}

// handleClientQuery 处理客户端查询消息
func (ch *ConnectionHandler) handleClientQuery(data []string, rawLine []byte) *Result {
	//	查询飞行计划
	//	$CQ ZYSH_CTR SERVER FP  CPA421
	//  [0] [  1   ] [  2 ] [3] [  4 ]
	//
	//	修改飞行计划
	//	$CQ ZYSH_CTR @94835 FA  CPA421 31100
	//	[0] [  1   ] [  2 ] [3] [  4 ] [ 5 ]
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	targetStation := data[1]
	if targetStation == "SERVER" {
		subQuery := data[2]
		// 查询指定机组的飞行计划
		if subQuery == "FP" {
			targetCallsign := data[3]
			client, ok := ch.clientManager.GetClient(targetCallsign)
			if !ok || client.FlightPlan() == nil {
				return ResultError(NoFlightPlan, false, ch.client.Callsign(), nil)
			}
			ch.client.SendLine([]byte(ch.flightPlanOperation.ToString(client.FlightPlan(), data[0])))
		}
	}
	// 如果发送目标是一个频率
	if strings.HasPrefix(targetStation, "@") {
		err := ch.sendFrequencyMessage(targetStation, rawLine)
		if err != nil {
			return err
		}
	} else {
		_ = ch.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

// handleClientResponse 处理客户端回复消息
func (ch *ConnectionHandler) handleClientResponse(data []string, rawLine []byte) *Result {
	//	$CR ZSHA_CTR ZSSS_APP CAPS ATCINFO=1 SECPOS=1 MODELDESC=1 ONGOINGCOORD=1 NEWINFO=1 TEAMSPEAK=1 ICAOEQ=1
	//  [0] [   1  ] [   2  ] [ 3] [   4   ] [  5   ] [    6    ] [     7      ] [   8   ] [     9   ] [  10  ]
	//	$CR ZSHA_CTR SERVER ATIS  T  ZSHA_CTR Shanghai Control
	//	[0] [   1  ] [  2 ] [ 3] [4] [           5           ]
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	commandLength := len(data)
	targetStation := data[1]
	if targetStation == "SERVER" {
		subQuery := data[2]
		if subQuery == "ATIS" && commandLength >= 5 && data[3] == "T" {
			ch.client.AddAtcAtisInfo(data[4])
		}
	}
	if strings.HasPrefix(targetStation, "@") {
		result := ch.sendFrequencyMessage(targetStation, rawLine)
		if result != nil {
			return result
		}
	} else {
		_ = ch.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

func (ch *ConnectionHandler) handleMessage(data []string, rawLine []byte) *Result {
	// #TM ZSHA_CTR ZSSS_APP 111
	// [0] [   1  ] [   2  ] [3]
	targetStation := data[1]
	if strings.HasPrefix(targetStation, "@") {
		result := ch.sendFrequencyMessage(targetStation, rawLine)
		if result != nil {
			return result
		}
	} else if strings.HasPrefix(targetStation, "*") {
		// 广播消息
		if targetStation == string(AllSup) {
			go ch.clientManager.BroadcastMessage(rawLine, ch.client, BroadcastToSup)
		}
	} else {
		_ = ch.clientManager.SendMessageTo(targetStation, rawLine)
	}
	return ResultSuccess()
}

func (ch *ConnectionHandler) handlePlan(data []string, rawLine []byte) *Result {
	// $FP CPA421 SERVER  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  1    18   2    26  ZYCC
	// [0] [  1 ] [  2 ] [3] [  4   ] [5] [ 6] [ 7] [8] [ 9 ] [10] [11] [12] [13] [14] [15]
	// /V/ SEL/AHFL VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [    16    ] [                      17                       ]
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if ch.client.IsAtc() {
		return ResultError(Syntax, false, ch.client.Callsign(), fmt.Errorf("atc can not submit fligth plan"))
	}
	if err := ch.client.UpsertFlightPlan(data); err != nil {
		return ResultError(Syntax, false, ch.client.Callsign(), err)
	}
	if !ch.client.FlightPlan().Locked {
		go ch.clientManager.BroadcastMessage(rawLine, ch.client, CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	}
	return ResultSuccess()
}

func (ch *ConnectionHandler) handleAtcEditPlan(data []string, _ []byte) *Result {
	// $AM ZYSH_CTR SERVER CPA421  I  H/A320/L 474 ZYTL 1115  0  FL371 ZYHB  11  8     22   6   ZYCC
	// [0] [   1  ] [  2 ] [  3 ] [4] [   5  ] [6] [ 7] [ 8] [9] [ 10] [11] [12] [13] [14] [15] [16]
	// /V/ SEL/AHFL CHI19D/28 VENOS A588 NULRA W206 MAGBI W656 ISLUK W629 LARUN
	// [     17   ] [                             18                          ]
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if !ch.client.IsAtc() {
		return ResultError(Syntax, false, ch.client.Callsign(), fmt.Errorf("only act can edit flight plan"))
	}
	if !ch.client.CheckFacility(AllowAtcFacility) {
		return ResultError(Syntax, false, ch.client.Callsign(), fmt.Errorf("%s facility not allowed to edit plan", ch.client.Facility().String()))
	}
	targetCallsign := data[2]
	client, ok := ch.clientManager.GetClient(targetCallsign)
	if !ok {
		return ResultError(SourceCallsignInvalid, false, ch.client.Callsign(), fmt.Errorf("%s not exists", targetCallsign))
	}
	if client.FlightPlan == nil {
		return ResultError(NoFlightPlan, false, ch.client.Callsign(), fmt.Errorf("%s do not have filght plan", ch.client.Callsign()))
	}
	client.FlightPlan().Locked = true
	if err := ch.flightPlanOperation.UpdateFlightPlan(client.FlightPlan(), client.Callsign(), data[1:], true); err != nil {
		return ResultError(Syntax, false, ch.client.Callsign(), err)
	}
	go ch.clientManager.BroadcastMessage([]byte(ch.flightPlanOperation.ToString(client.FlightPlan(), string(AllATC))),
		ch.client, CombineBroadcastFilter(BroadcastToAtc, BroadcastToClientInRange))
	return ResultSuccess()
}

func (ch *ConnectionHandler) handleKillClient(data []string, _ []byte) *Result {
	// $!! ZSHA_CTR CPA421 test
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	if !(ch.client.IsAtc() && ch.client.CheckRating(AllowKillRating)) {
		return ResultError(Syntax, false, ch.client.Callsign(), fmt.Errorf("%s facility not allowed to kill client", ch.client.Facility().String()))
	}
	targetStation := data[1]
	client, ok := ch.clientManager.GetClient(targetStation)
	if !ok {
		return ResultError(NoCallsignFound, false, ch.client.Callsign(), fmt.Errorf("%s not exists", targetStation))
	}
	client.MarkedDisconnect(false)
	return ResultSuccess()
}

func (ch *ConnectionHandler) handleRequest(data []string, rawLine []byte) *Result {
	// $HO ZSSS_APP ZSHA_CTR CES2352
	// [0] [  1   ] [   2  ] [  3  ]
	// #PC ZSHA_CTR ZSSS_APP CCP HC CES2352
	// $HA ZSHA_CTR ZSSS_APP CES2352
	targetStation := data[1]
	_ = ch.clientManager.SendMessageTo(targetStation, rawLine)
	return ResultSuccess()
}

func (ch *ConnectionHandler) removeClient(_ []string, _ []byte) *Result {
	// #DA ZGGG_CTR SERVER
	c.InfoF("[%s] Offline", ch.client.Callsign())
	return ResultSuccess()
}

func (ch *ConnectionHandler) handleSquawkBox(data []string, rawLine []byte) *Result {
	//# SB ZSHA_CTR CES7199 FSIPIR 0 ZZZ C172 3.73453 1.32984 174.00000 3.3028DC0.98CF9A41 L1P Cessna Skyhawk 172SP
	if ch.client == nil {
		return ResultError(Syntax, false, "", fmt.Errorf("client not register"))
	}
	targetStation := data[1]
	client, ok := ch.clientManager.GetClient(targetStation)
	if !ok {
		return ResultError(SourceCallsignInvalid, false, ch.client.Callsign(), fmt.Errorf("%s not exists", targetStation))
	}
	client.SendLine(rawLine)
	return ResultSuccess()
}

func (ch *ConnectionHandler) handleCommand(commandType ClientCommand, data []string, rawLine []byte) *Result {
	var result = ResultSuccess()
	if requirement, ok := CommandRequirements[commandType]; ok {
		if err, ok := ch.checkPacketLength(data, requirement); !ok {
			return err
		}
	}
	switch commandType {
	case AddAtc:
		result = ch.handleAddAtc(data, rawLine)
	case AddPilot:
		result = ch.handleAddPilot(data, rawLine)
	case RequestHandoff, AcceptHandoff, ProController:
		result = ch.handleRequest(data, rawLine)
	case PilotPosition:
		result = ch.handlePilotPosUpdate(data, rawLine)
	case Plan:
		result = ch.handlePlan(data, rawLine)
	case AtcEditPlan:
		result = ch.handleAtcEditPlan(data, rawLine)
	case KillClient:
		result = ch.handleKillClient(data, rawLine)
	case AtcPosition:
		result = ch.handleAtcPosUpdate(data, rawLine)
	case AtcSubVisPoint:
		result = ch.handleAtcVisPointUpdate(data, rawLine)
	case Message:
		result = ch.handleMessage(data, rawLine)
	case ClientQuery:
		result = ch.handleClientQuery(data, rawLine)
	case ClientResponse:
		result = ch.handleClientResponse(data, rawLine)
	case RemoveAtc, RemovePilot:
		ch.removeClient(data, rawLine)
	case SquawkBox:
		result = ch.handleSquawkBox(data, rawLine)
	default:
		result = ResultSuccess()
	}
	return result
}
