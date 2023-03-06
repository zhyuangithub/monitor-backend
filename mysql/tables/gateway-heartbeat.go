package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

type GatewayHeartBeat struct {
	GatewayId int64  `json:"gatewayId"`
	Ipv4      string `json:"ipv4"`
	Ipv6      string `json:"ipv6"`
	Timestamp string `json:"timestamp"`
}

func (g *GatewayHeartBeat) getTableName() string {
	return utils.GetColByAction(utils.GatewayHeartBeat)
}
func (g *GatewayHeartBeat) GenerateCreateTableStr() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		gatewayId INT NOT NULL, ipv4 TEXT, ipv6 TEXT, timestamp TIMESTAMP)`, g.getTableName())
}
func (g *GatewayHeartBeat) GenerateInsertStr(data map[string]interface{}) string {
	res, _ := utils.MapToStruct(data, &GatewayHeartBeat{})
	val := res.(*GatewayHeartBeat)
	return fmt.Sprintf(`insert into %s(gatewayId,ipv4,ipv6,timestamp) values(%d,"%s","%s","%s");`, g.getTableName(),
		val.GatewayId, val.Ipv4, val.Ipv6, val.Timestamp)
}
