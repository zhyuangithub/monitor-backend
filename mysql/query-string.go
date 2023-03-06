package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

func generateQueryBatchNodesStr(isCount bool, name string, status string) string {
	nodeStates := utils.GetColByAction(utils.NodeStates)
	nodeName := utils.GetColByAction(utils.NodeName)
	sqlStr := ""
	countStr := "*"
	if isCount {
		countStr = "count(*)"
	}
	if name != "" && status != "" {
		sqlStr = fmt.Sprintf(`select %s from %s s, %s t`, countStr, nodeName, nodeStates)
		sqlStr = fmt.Sprintf(`%s where s.nodeName like '%%%s%%' and s.monitorId = t.state->'$.monitorId'`, sqlStr, name)
		if status == "online" {
			sqlStr = fmt.Sprintf(`%s and t.state->'$.thermal' is not null`, sqlStr)
		} else if status == "offline" {
			sqlStr = fmt.Sprintf(`%s and t.state->'$.thermal' is null`, sqlStr)
		}
	} else if name == "" && status == "" {
		sqlStr = fmt.Sprintf(`select %s from %s`, countStr, nodeStates)
	} else {
		sqlStr = fmt.Sprintf(`select %s from %s`, countStr, nodeStates)
		if status == "online" {
			sqlStr = fmt.Sprintf(`%s where json_extract(state,'$.thermal') is not null`, sqlStr)
		} else if status == "offline" {
			sqlStr = fmt.Sprintf(`%s where json_extract(state,'$.thermal') is null`, sqlStr)
		} else {
			sqlStr = fmt.Sprintf(`select %s from %s s, %s t`, countStr, nodeName, nodeStates)
			sqlStr = fmt.Sprintf(`%s where s.nodeName like '%%%s%%' and s.monitorId = t.state->'$.monitorId'`, sqlStr, name)
		}
	}
	return sqlStr
}
