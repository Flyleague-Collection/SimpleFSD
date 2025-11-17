// Package pdu
package pdu

import "github.com/half-nothing/simple-fsd/src/interfaces/enum"

type QueryType *enum.Enum[string, string]

var (
	QueryTypeInvalid               QueryType = enum.New("", "")
	QueryTypeIsValidATC            QueryType = enum.New("ATC", "查询是否是合法管制员")
	QueryTypeCapabilities          QueryType = enum.New("CAPS", "查询客户端能力")
	QueryTypeCOM1Freq              QueryType = enum.New("C?", "查询客户端Freq1频率")
	QueryTypeRealName              QueryType = enum.New("RN", "查询客户端真实姓名")
	QueryTypeServer                QueryType = enum.New("SV", "查询服务器信息")
	QueryTypeATIS                  QueryType = enum.New("ATIS", "查询ATIS信息")
	QueryTypePublicIP              QueryType = enum.New("IP", "查询客户端公网IP")
	QueryTypeINF                   QueryType = enum.New("INF", "查询客户端信息")
	QueryTypeFlightPlan            QueryType = enum.New("FP", "查询客户端计划")
	QueryTypeIPC                   QueryType = enum.New("IPC", "查询客户端IPC信息")
	QueryTypeRequestRelief         QueryType = enum.New("BY", "管制员暂离请求")
	QueryTypeCancelRequestRelief   QueryType = enum.New("HI", "管制员结束暂离")
	QueryTypeRequestHelp           QueryType = enum.New("HLP", "请求帮助")
	QueryTypeCancelRequestHelp     QueryType = enum.New("NOHLP", "取消请求帮助")
	QueryTypeWhoHas                QueryType = enum.New("WH", "查询谁有牌子")
	QueryTypeInitiateTrack         QueryType = enum.New("IT", "接牌")
	QueryTypeAcceptHandoff         QueryType = enum.New("HT", "接受移交")
	QueryTypeDropTrack             QueryType = enum.New("DR", "丢牌子")
	QueryTypeSetFinalAltitude      QueryType = enum.New("FA", "设置巡航高度")
	QueryTypeSetTempAltitude       QueryType = enum.New("TA", "设置临时高度")
	QueryTypeSetSquawkCode         QueryType = enum.New("BC", "设置代码")
	QueryTypeSetScratchpad         QueryType = enum.New("SC", "设置备注")
	QueryTypeSetVoiceType          QueryType = enum.New("VT", "设置交流类型")
	QueryTypeAircraftConfiguration QueryType = enum.New("ACC", "查询客户端状态")
	QueryTypeNewInfo               QueryType = enum.New("NEWINFO", "管制员信息更新")
	QueryTypeNewATIS               QueryType = enum.New("NEWATIS", "ATIS更新")
	QueryTypeEstimate              QueryType = enum.New("EST", "预计时间")
	QueryTypeSetGlobalData         QueryType = enum.New("GD", "设置全局数据")
)

var QueryTypes = enum.NewManager(
	QueryTypeInvalid,
	QueryTypeIsValidATC,
	QueryTypeCapabilities,
	QueryTypeCOM1Freq,
	QueryTypeRealName,
	QueryTypeServer,
	QueryTypeATIS,
	QueryTypePublicIP,
	QueryTypeINF,
	QueryTypeFlightPlan,
	QueryTypeIPC,
	QueryTypeRequestRelief,
	QueryTypeCancelRequestRelief,
	QueryTypeRequestHelp,
	QueryTypeCancelRequestHelp,
	QueryTypeWhoHas,
	QueryTypeInitiateTrack,
	QueryTypeAcceptHandoff,
	QueryTypeDropTrack,
	QueryTypeSetFinalAltitude,
	QueryTypeSetTempAltitude,
	QueryTypeSetSquawkCode,
	QueryTypeSetScratchpad,
	QueryTypeSetVoiceType,
	QueryTypeAircraftConfiguration,
	QueryTypeNewInfo,
	QueryTypeNewATIS,
	QueryTypeEstimate,
	QueryTypeSetGlobalData,
)
