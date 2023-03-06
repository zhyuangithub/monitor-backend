package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

type MonitorState struct {
	//Id                 int64   `json:"id"`
	MonitorId          int64   `json:"monitorId"`
	ViaGatewayId       int64   `json:"viaGatewayId"`
	Rssi               float64 `json:"rssi"`
	Snr                float64 `json:"snr"`
	Vbat               float64 `json:"vbat"`
	AmbientTemperature float64 `json:"ambientTemperature"`
	ObjectTemperature  float64 `json:"objectTemperature"`
	Thermal            bool    `json:"thermal"`
	Voltage            float64 `json:"voltage"`
	Current            float64 `json:"current"`
	Power              float64 `json:"power"`
	Kwh                float64 `json:"kwh"`
	Timestamp          string  `json:"timestamp"`
}

// 只存最新的一条
// last updated at需要转换时区，存入Timestamp
func (m *MonitorState) getTableName() string {
	return utils.GetColByAction(utils.MonitorState)
}
func (m *MonitorState) GenerateCreateTableStr() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		monitorId INT NOT NULL, viaGatewayId INT, rssi FLOAT, snr FLOAT, vbat FLOAT, 
		ambientTemperature FLOAT, objectTemperature FLOAT, thermal TINYINT, voltage FLOAT, current FLOAT,
		power FLOAT, kwh FLOAT, timestamp TIMESTAMP)`, m.getTableName())
}

func (m *MonitorState) GenerateInsertStr(data map[string]interface{}) string {
	res, _ := utils.MapToStruct(data, &MonitorState{})
	val := res.(*MonitorState)
	themal := 0
	if val.Thermal {
		themal = 1
	}
	return fmt.Sprintf(`insert into %s(monitorId,viaGatewayId,rssi,snr,vbat,ambientTemperature,objectTemperature,thermal,voltage,current,power,kwh,timestamp) values(%d,%d,%f,%f,%f,%f,%f,%d,%f,%f,%f,%f,"%s");`,
		m.getTableName(), val.MonitorId, val.ViaGatewayId, val.Rssi, val.Snr, val.Vbat, val.AmbientTemperature, val.ObjectTemperature, themal, val.Voltage, val.Current, val.Power, val.Kwh, val.Timestamp)
}
func (m *MonitorState) GenerateUpdateStr(data map[string]interface{}) string {
	res, _ := utils.MapToStruct(data, &MonitorState{})
	val := res.(*MonitorState)
	themal := 0
	if val.Thermal {
		themal = 1
	}
	return fmt.Sprintf(`update %s set viaGatewayId=%d,rssi=%f,snr=%f,vbat=%f,ambientTemperature=%f,objectTemperature=%f,thermal=%d,voltage=%f,current=%f,power=%f,kwh=%f,timestamp='%s' where monitorId=%d;`,
		m.getTableName(), val.ViaGatewayId, val.Rssi, val.Snr, val.Vbat, val.AmbientTemperature, val.ObjectTemperature, themal, val.Voltage, val.Current, val.Power, val.Kwh, val.Timestamp, val.MonitorId)
}

/*
	func (m *MonitorState) GenerateQueryStr() string {
		//select count(*) from monitor_state where id=3;
		return fmt.Sprintf(`insert into %s(monitorId,viaGatewayId,rssi,snr,vbat,ambientTemperature,objectTemperature,thermal,voltage,current,power,kwh,timestamp) values(%d,%d,%f,%f,%f,%f,%f,%d,%f,%f,%f,%f,"%s");`,
			m.getTableName(), val.MonitorId, val.ViaGatewayId, val.Rssi, val.Snr, val.Vbat, val.AmbientTemperature, val.ObjectTemperature, themal, val.Voltage, val.Current, val.Power, val.Kwh, val.Timestamp)
	}
*/

func (m *MonitorState) GenerateCheckIdStr(monitorId int64) string {
	return fmt.Sprintf(`select count(*) from %s where monitorId=%d;`, m.getTableName(), monitorId)
}
func (m *MonitorState) GenerateDeleteStr(monitorId int64) string {
	return fmt.Sprintf(`delete from %s where monitorId=%d;`, m.getTableName(), monitorId)
}
