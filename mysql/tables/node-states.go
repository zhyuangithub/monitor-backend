package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

type NodeStates struct {
	State string `json:"state"`
}

func (n *NodeStates) getTableName() string {
	return utils.GetColByAction(utils.NodeStates)
}
func (n *NodeStates) GenerateCreateTableStr() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		state json NOT NULL)`, n.getTableName())
}
func (n *NodeStates) GenerateInsertStr(data map[string]interface{}) string {
	jsonStr, _ := utils.MapToJson(data)
	return fmt.Sprintf(`insert into %s (state) values('%s');`, n.getTableName(), jsonStr)
}
func (n *NodeStates) GenerateCheckIdStr(monitorId int64) string {
	return fmt.Sprintf(`select count(*) from %s where state->'$.monitorId' = %d;`, n.getTableName(), monitorId)
}
func (n *NodeStates) GenerateUpdateStr(data map[string]interface{}) string {
	jsonStr, _ := utils.MapToJson(data)
	monitorId := int64(data["monitorId"].(float64))
	return fmt.Sprintf(`update %s set state='%s' where state->'$.monitorId'=%d;`, n.getTableName(), jsonStr, monitorId)
}
func (n *NodeStates) GenerateDeleteStr(monitorId int64) string {
	return fmt.Sprintf(`delete from %s where state->'$.monitorId' = %d;`, n.getTableName(), monitorId)
}
