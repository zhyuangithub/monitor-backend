package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

type GatewayState struct {
	GatewayId       int64  `json:"gatewayId"`
	FirmwareVersion string `json:"firmwareVersion"`
	ConnectionType  string `json:"connectionType"`
	Status          string `json:"status"`
	Ipv4            string `json:"ipv4"`
	Ipv6            string `json:"ipv6"`
	Timestamp       string `json:"timestamp"`
}

func (g *GatewayState) getTableName() string {
	return utils.GetColByAction(utils.GatewayState)
}
func (g *GatewayState) GenerateCreateTableStr() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		gatewayId INT NOT NULL, firmwareVersion TEXT, connectionType TEXT, status TEXT, ipv4 TEXT, ipv6 TEXT, timestamp TIMESTAMP)`, g.getTableName())
}
func (g *GatewayState) GenerateInsertStr(data map[string]interface{}) string {
	res, _ := utils.MapToStruct(data, &GatewayState{})
	val := res.(*GatewayState)
	return fmt.Sprintf(`insert into %s(gatewayId,firmwareVersion,connectionType,status,ipv4,ipv6,timestamp) values(%d,"%s","%s","%s","%s","%s","%s");`, g.getTableName(),
		val.GatewayId, val.FirmwareVersion, val.ConnectionType, val.Status, val.Ipv4, val.Ipv6, val.Timestamp)
}
