package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

type NodeName struct {
	MonitorId int64  `json:"monitorId"`
	NodeName  string `json:"nodeName"`
}

func (n *NodeName) getTableName() string {
	return utils.GetColByAction(utils.NodeName)
}
func (n *NodeName) GenerateCreateTableStr() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		monitorId INT NOT NULL, nodeName TEXT) DEFAULT CHARSET=utf8;`, n.getTableName())
}
func (n *NodeName) GenerateInsertStr(data map[string]interface{}) string {
	res, _ := utils.MapToStruct(data, &NodeName{})
	val := res.(*NodeName)
	return fmt.Sprintf(`insert into %s(monitorId,nodeName) values(%d,"%s");`, n.getTableName(),
		val.MonitorId, val.NodeName)
}
func (n *NodeName) GenerateCheckIdStr(monitorId int64) string {
	return fmt.Sprintf(`select count(*) from %s where monitorId=%d;`, n.getTableName(), monitorId)
}
func (n *NodeName) GenerateDeleteStr(monitorId int64) string {
	return fmt.Sprintf(`delete from %s where monitorId=%d;`, n.getTableName(), monitorId)
}
func (n *NodeName) GenerateUpdateStr(data map[string]interface{}) string {
	res, _ := utils.MapToStruct(data, &NodeName{})
	val := res.(*NodeName)
	return fmt.Sprintf(`update %s set nodeName='%s' where monitorId=%d;`,
		n.getTableName(), val.NodeName, val.MonitorId)
}
