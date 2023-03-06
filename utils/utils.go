package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	Version          = "1.2"
	GatewayState     = "GatewayState"     //网关状态发生改变，例如：上线、离线
	GatewayHeartBeat = "GatewayHeartBeat" //网关心跳
	MonitorState     = "MonitorState"     //节点数据更新
	Dismount         = "Dismount"         //节点防拆开关触发
	NFCCardRead      = "NFCCardRead"      //NFC刷卡触发
	MonitorSleeping  = "MonitorSleeping"  //节点进入睡眠状态
	NodeName         = "NodeName"         //节点进入睡眠状态
	Events           = "Events"
	NodeStates       = "NodeStates"
)

func GetColByAction(action string) string {
	switch action {
	case GatewayState:
		return "gateway_state"
	case GatewayHeartBeat:
		return "gateway_heartbeat"
	case MonitorState:
		return "monitor_state"
	case Dismount:
		return "dismount"
	case NFCCardRead:
		return "nfc_card_read"
	case MonitorSleeping:
		return "monitor_sleeping"
	case NodeName:
		return "node_name"
	case Events:
		return "events"
	case NodeStates:
		return "node_states"
	default:
		return ""
	}

}
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v", err)
		return nil, err
	}
	return m, nil
}
func MapToJson(data map[string]interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return "", err
	}
	return string(b), nil
}
func TimeStrToTimestamp(timeStr string) (int64, error) {
	res, err := TimeStrToTime(timeStr)
	return res.Unix(), err
}
func TimeStrToTime(timeStr string) (time.Time, error) {
	layoutStr := "2006-01-02 15:04:05"
	local, _ := time.LoadLocation("Asia/Shanghai")
	parsed, err := time.ParseInLocation(layoutStr, timeStr, local)
	if err != nil {
		parsed, err = IsoTimeStrToTime(timeStr)
	}
	return parsed, err
}

func TimeStampToLocalTimeStr(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	layoutStr := "2006-01-02 15:04:05"
	local, _ := time.LoadLocation("Asia/Shanghai")
	return tm.In(local).Format(layoutStr)
}
func IsoTimeStrToTime(timeStr string) (time.Time, error) {
	layoutStr := "2006-01-02T15:04:05.000Z"
	parsed, err := time.Parse(layoutStr, timeStr)
	return parsed, err
}
func MapToStruct(val any, res any) (any, error) {
	arr, err := json.Marshal(val)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(arr, &res)
	return res, err
}
func StructToMap(val any) (map[string]interface{}, error) {
	data, err := json.Marshal(&val)
	res := make(map[string]interface{})
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	return res, err
}
