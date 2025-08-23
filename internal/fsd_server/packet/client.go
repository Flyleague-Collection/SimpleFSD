package packet

import (
	"bytes"
	"errors"
	"fmt"
	c "github.com/half-nothing/simple-fsd/internal/config"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/fsd"
	"github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"github.com/half-nothing/simple-fsd/internal/utils"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	isAtc               bool
	callsign            string
	rating              Rating
	facility            Facility
	user                *operation.User
	protocol            int
	realName            string
	socket              ConnectionHandlerInterface
	position            [4]Position
	simType             int
	transponder         string
	altitude            int
	groundSpeed         int
	frequency           int
	pbh                 uint32
	visualRange         float64
	flightPlan          *operation.FlightPlan
	atisInfo            []string
	paths               []*PilotPath
	history             *operation.History
	clientManager       ClientManagerInterface
	disconnect          atomic.Bool
	motdBytes           []byte
	reconnectTimer      *time.Timer
	lock                sync.RWMutex
	config              *c.Config
	userOperation       operation.UserOperationInterface
	flightPlanOperation operation.FlightPlanOperationInterface
	historyOperation    operation.HistoryOperationInterface
	pathTrigger         *utils.OverflowTrigger
}

func (cm *ClientManager) NewClient(callsign string, rating Rating, protocol int, realName string, socket ConnectionHandlerInterface, isAtc bool) ClientInterface {
	socket.SetCallsign(callsign)
	var flightPlan *operation.FlightPlan = nil
	flightPlanOperation := cm.applicationContent.FlightPlanOperation()
	userOperation := cm.applicationContent.UserOperation()
	historyOperation := cm.applicationContent.HistoryOperation()
	if !isAtc && !cm.config.Server.General.SimulatorServer {
		var err error
		flightPlan, err = flightPlanOperation.GetFlightPlanByCid(socket.User().Cid)
		if errors.Is(err, operation.ErrFlightPlanNotFound) {
			c.WarnF("No flight plan found for %s(%d)", callsign, socket.User().Cid)
		} else if err != nil {
			c.WarnF("Fail to get flight plan for %s(%d): %v", callsign, socket.User().Cid, err)
		}
	}
	client := &Client{
		isAtc:               isAtc,
		callsign:            callsign,
		rating:              rating,
		facility:            0,
		user:                socket.User(),
		protocol:            protocol,
		realName:            realName,
		socket:              socket,
		position:            [4]Position{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
		simType:             0,
		transponder:         "2000",
		altitude:            0,
		groundSpeed:         0,
		frequency:           99998,
		visualRange:         40,
		flightPlan:          flightPlan,
		atisInfo:            make([]string, 0, 4),
		paths:               make([]*PilotPath, 0),
		history:             historyOperation.NewHistory(socket.User().Cid, callsign, isAtc),
		motdBytes:           nil,
		clientManager:       cm,
		disconnect:          atomic.Bool{},
		reconnectTimer:      nil,
		lock:                sync.RWMutex{},
		config:              cm.config,
		userOperation:       userOperation,
		flightPlanOperation: flightPlanOperation,
		historyOperation:    historyOperation,
	}
	client.pathTrigger = utils.NewOverflowTrigger(cm.config.Server.FSDServer.PosUpdatePoints, client.recordPathPoint)
	return client
}

func (client *Client) recordPathPoint() {
	client.paths = append(client.paths, &PilotPath{
		Latitude:  client.position[0].Latitude,
		Longitude: client.position[0].Longitude,
		Altitude:  client.altitude,
	})
}

func (client *Client) Disconnected() bool {
	return client.disconnect.Load()
}

func (client *Client) Delete() {
	if client.disconnect.Load() {
		client.lock.Lock()
		defer client.lock.Unlock()
		c.InfoF("[%s](%s) client session deleted", client.socket.ConnId(), client.callsign)

		if client.reconnectTimer != nil {
			client.reconnectTimer.Stop()
			client.reconnectTimer = nil
		}

		if client.isAtc || !client.config.Server.General.SimulatorServer {
			if err := client.historyOperation.EndRecordAndSaveHistory(client.history); err != nil {
				c.ErrorF("[%s](%s) Failed to end history: %v", client.socket.ConnId(), client.callsign, err)
			}
		}

		if client.isAtc {
			if err := client.userOperation.UpdateUserAtcTime(client.user, client.history.OnlineTime); err != nil {
				c.ErrorF("[%s](%s) Failed to add ATC time: %v", client.socket.ConnId(), client.callsign, err)
			}
		} else if !client.config.Server.General.SimulatorServer {
			// 如果不是模拟机服务器, 则写入机组连线时长
			if err := client.userOperation.UpdateUserPilotTime(client.user, client.history.OnlineTime); err != nil {
				c.ErrorF("[%s](%s) Failed to add pilot time: %v", client.socket.ConnId(), client.callsign, err)
			}
		}

		if !client.clientManager.DeleteClient(client.callsign) {
			c.ErrorF("[%s](%s) Failed to delete from client manager", client.socket.ConnId(), client.callsign)
		}
	}
}

func (client *Client) Reconnect(socket ConnectionHandlerInterface) bool {
	client.lock.Lock()
	defer client.lock.Unlock()

	if !client.disconnect.Load() {
		return false
	}

	c.InfoF("[%s](%s) client reconnected", client.socket.ConnId, client.callsign)

	if client.reconnectTimer != nil {
		client.reconnectTimer.Stop()
		client.reconnectTimer = nil
	}

	client.ClearAtcAtisInfo()
	client.disconnect.Store(false)
	client.socket = socket
	socket.SetCallsign(client.callsign)
	return true
}

func (client *Client) MarkedDisconnect(immediate bool) {
	client.lock.Lock()
	defer func() {
		client.lock.Unlock()
		if immediate {
			client.Delete()
		}
	}()

	if !client.disconnect.CompareAndSwap(false, true) {
		return
	}

	// 关闭连接
	if client.socket.Conn() != nil {
		_ = client.socket.Conn().Close()
	}

	// 取消之前的定时器
	if client.reconnectTimer != nil {
		client.reconnectTimer.Stop()
	}

	if immediate {
		return
	}

	client.reconnectTimer = time.AfterFunc(client.config.Server.FSDServer.SessionCleanDuration, client.Delete)
	c.InfoF("[%s](%s) client disconnected, reconnect window: %v", client.socket.ConnId(),
		client.callsign, client.config.Server.FSDServer.SessionCleanDuration)
}

func (client *Client) UpsertFlightPlan(flightPlanData []string) error {
	if client.flightPlan == nil {
		flightPlan, err := client.flightPlanOperation.UpsertFlightPlan(client.user, client.callsign, flightPlanData)
		if err != nil {
			return err
		}
		client.flightPlan = flightPlan
		return nil
	}
	// 如果是模拟机服务器, 只创建就行
	if client.config.Server.General.SimulatorServer {
		return nil
	}
	if client.flightPlan.Locked {
		departureAirport := flightPlanData[5]
		arrivalAirport := flightPlanData[9]
		if client.flightPlan.DepartureAirport != departureAirport || client.flightPlan.ArrivalAirport != arrivalAirport {
			client.flightPlan.Locked = false
		}
	}
	err := client.flightPlanOperation.UpdateFlightPlan(client.flightPlan, flightPlanData, false)
	return err
}

func (client *Client) SetPosition(index int, lat float64, lon float64) error {
	if index >= 4 {
		return errors.New("position index out of range")
	}
	client.position[index].Latitude = lat
	client.position[index].Longitude = lon
	return nil
}

func (client *Client) UpdatePilotPos(transponder int, lat float64, lon float64, alt int, groundSpeed int, pbh uint32) {
	_ = client.SetPosition(0, lat, lon)
	client.transponder = fmt.Sprintf("%04d", transponder)
	client.altitude = alt
	client.groundSpeed = groundSpeed
	client.pbh = pbh
	go client.pathTrigger.Tick()
}

func (client *Client) UpdateAtcPos(frequency int, facility Facility, visualRange float64, lat float64, lon float64) {
	_ = client.SetPosition(0, lat, lon)
	client.frequency = frequency
	client.facility = facility
	client.visualRange = visualRange
}

func (client *Client) UpdateAtcVisPoint(visIndex int, lat float64, lon float64) error {
	if visIndex < 0 || visIndex > 2 {
		return errors.New("visIndex out of range [0,2]")
	}
	return client.SetPosition(visIndex+1, lat, lon)
}

func (client *Client) ClearAtcAtisInfo() {
	client.atisInfo = client.atisInfo[:0]
}

func (client *Client) AddAtcAtisInfo(atisInfo string) {
	client.atisInfo = append(client.atisInfo, atisInfo)
}

func (client *Client) SendError(result *Result) {
	if result.Success {
		return
	}

	packet := makePacket(Error, "fsd_server", client.callsign, fmt.Sprintf("%03d", result.Errno.Index()), result.Env, result.Errno.String())
	client.SendLine(packet)

	if result.Fatal {
		client.socket.SetDisconnected(true)
		client.disconnect.Store(true)
		time.AfterFunc(500*time.Millisecond, func() {
			if !client.clientManager.DeleteClient(client.callsign) {
				c.ErrorF("[%s](%s) Failed to delete from client manager", client.socket.ConnId(), client.callsign)
			}
		})
	}
}

func (client *Client) SendLineWithoutLog(line []byte) {
	client.lock.RLock()
	defer client.lock.RUnlock()

	if client.disconnect.Load() {
		c.WarnF("[%s](%s) Attempted send to disconnected client", client.socket.ConnId(), client.callsign)
		return
	}

	if !bytes.HasSuffix(line, splitSign) {
		line = append(line, splitSign...)
	}

	if _, err := client.socket.Conn().Write(line); err != nil {
		c.ErrorF("[%s](%s) Failed to send data: %v", client.socket.ConnId(), client.callsign, err)
	}
}

func (client *Client) SendLine(line []byte) {
	client.lock.RLock()
	defer client.lock.RUnlock()

	if client.disconnect.Load() {
		c.DebugF("[%s](%s) Attempted send to disconnected client", client.socket.ConnId(), client.callsign)
		return
	}

	if !bytes.HasSuffix(line, splitSign) {
		c.DebugF("[%s](%s) <- %s", client.socket.ConnId(), client.callsign, line)
		line = append(line, splitSign...)
	} else {
		c.DebugF("[%s](%s) <- %s", client.socket.ConnId(), client.callsign, line[:len(line)-splitSignLen])
	}

	if _, err := client.socket.Conn().Write(line); err != nil {
		c.WarnF("[%s](%s) Failed to send data: %v", client.socket.ConnId(), client.callsign, err)
	}
}

func (client *Client) SendMotd() {
	if client.motdBytes != nil {
		client.SendLine(client.motdBytes)
		return
	}

	data := make([][]byte, 0, len(client.config.Server.FSDServer.Motd)+1)
	data = append(data, []byte(fmt.Sprintf("%sserver:%s:Welcome to use %s v%s\r\n",
		Message, client.callsign, client.config.Server.FSDServer.FSDName, c.AppVersion.String())))

	for _, message := range client.config.Server.FSDServer.Motd {
		data = append(data, makePacket(Message, "fsd_server", client.callsign, message))
	}

	buffer := bytes.Buffer{}
	for _, msg := range data {
		buffer.Write(msg)
	}
	client.motdBytes = buffer.Bytes()
	client.SendLine(client.motdBytes)
}

func (client *Client) CheckFacility(facility Facility) bool {
	return facility.CheckFacility(client.facility)
}

func (client *Client) CheckRating(rating []Rating) bool {
	return slices.Contains(rating, client.rating)
}

func (client *Client) IsAtc() bool { return client.isAtc }

func (client *Client) Callsign() string { return client.callsign }

func (client *Client) Rating() Rating { return client.rating }

func (client *Client) Facility() Facility { return client.facility }

func (client *Client) RealName() string { return client.realName }

func (client *Client) Position() [4]Position { return client.position }

func (client *Client) VisualRange() float64 { return client.visualRange }

func (client *Client) SetUser(user *operation.User) { client.user = user }

func (client *Client) SetSimType(simType int) { client.simType = simType }

func (client *Client) FlightPlan() *operation.FlightPlan { return client.flightPlan }

func (client *Client) User() *operation.User { return client.user }

func (client *Client) Frequency() int { return client.frequency }

func (client *Client) AtisInfo() []string { return client.atisInfo }

func (client *Client) History() *operation.History { return client.history }

func (client *Client) Transponder() string { return client.transponder }

func (client *Client) Altitude() int { return client.altitude }

func (client *Client) GroundSpeed() int { return client.groundSpeed }

func (client *Client) Heading() int {
	_, _, heading, _ := utils.UnpackPBH(client.pbh)
	return int(heading)
}

func (client *Client) Paths() []*PilotPath {
	return client.paths
}
