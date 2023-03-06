package mysql

import (
	"fmt"
	"monitor-backend/utils"
)

type Events struct {
	Event string `json:"event"`
}

func (e *Events) getTableName() string {
	return utils.GetColByAction(utils.Events)
}
func (e *Events) GenerateCreateTableStr() string {
	return fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, 
		event json NOT NULL)`, e.getTableName())
}
func (e *Events) GenerateInsertStr(data map[string]interface{}) string {
	jsonStr, _ := utils.MapToJson(data)
	return fmt.Sprintf(`insert into %s (event) values('%s');`, e.getTableName(), jsonStr)
}
