package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	tables "monitor-backend/mysql/tables"
	"monitor-backend/utils"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var instance *Mysql

type Mysql struct {
	db     *sql.DB
	tables map[string]interface{}
}

func MysqlInstance() *Mysql {
	if instance == nil {
		instance = new(Mysql)
		instance.init()
	}
	return instance
}
func (m *Mysql) init() {
	m.tables = make(map[string]interface{})
	url := fmt.Sprintf("%s/%s", os.Getenv("MYSQL_URL"), os.Getenv("MYSQL_DATABASE"))
	m.db, _ = sql.Open("mysql", url)
	err := m.db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
	m.initTables()
}
func (m *Mysql) initTables() {
	m.tables[utils.GatewayState] = &tables.GatewayState{}
	m.tables[utils.GatewayHeartBeat] = &tables.GatewayHeartBeat{}
	m.tables[utils.NodeStates] = &tables.NodeStates{}
	m.tables[utils.NodeName] = &tables.NodeName{}
	m.tables[utils.Events] = &tables.Events{}

	tableNameList := m.showTables()
	totalTables := []string{utils.GatewayState, utils.GatewayHeartBeat, utils.NodeStates,
		utils.NodeName, utils.Events}
	for _, data := range totalTables {
		tableName := utils.GetColByAction(data)
		if !checkTableExist(tableNameList, tableName) {
			table := m.tables[data].(tables.ITable)
			m.createTable(table.GenerateCreateTableStr())
		}
	}
}
func (m *Mysql) createTable(sqlStr string) {
	_, err := m.db.Exec(sqlStr)
	if err != nil {
		fmt.Println(err.Error()) //失败
	}
}
func (m *Mysql) InsertData(action string, data map[string]interface{}) error {
	if action == "" {
		return errors.New("invalid action!")
	}
	switch action {
	case utils.GatewayState, utils.GatewayHeartBeat, utils.NodeName:
		table := m.tables[action].(tables.ITable)
		sqlStr := table.GenerateInsertStr(data)
		return m.execSql(sqlStr)
	case utils.Dismount, utils.MonitorSleeping, utils.NFCCardRead:
		//timeStampStr转成东八区的
		table := m.tables[utils.Events].(tables.ITable)
		eventsData := make(map[string]interface{})
		eventsData["action"] = action
		eventsData["data"] = data
		sqlStr := table.GenerateInsertStr(eventsData)
		return m.execSql(sqlStr)
	default:
		return nil
	}

}
func (m *Mysql) execSql(sqlStr string) error {
	res, err := m.db.Exec(sqlStr)
	if err != nil {
		return err
	} else {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New("exec data failed!")
			//failed
		} else {
			return nil
		}
	}
}

func (m *Mysql) StoreNodeStates(data map[string]interface{}) error {
	//存在json数据里面，monitorstate或monitorsleep
	table := m.tables[utils.NodeStates].(*tables.NodeStates)
	monitorId := int64(data["monitorId"].(float64))
	sqlStr := table.GenerateCheckIdStr(monitorId)
	var count int
	m.db.QueryRow(sqlStr).Scan(&count)
	if count == 1 {
		//update
		sqlStr := table.GenerateUpdateStr(data)
		res, err := m.db.Exec(sqlStr)
		if err != nil {
			return err
		} else {
			rows, _ := res.RowsAffected()
			if rows == 0 {
				return errors.New("update data failed!")
				//failed
			} else {
				return nil
			}
		}
	} else if count == 0 {
		insertSql := table.GenerateInsertStr(data)
		return m.execSql(insertSql)
	} else {
		deleteSql := table.GenerateDeleteStr(monitorId)
		m.execSql(deleteSql)
		insertSql := table.GenerateInsertStr(data)
		return m.execSql(insertSql)
	}
}
func (m *Mysql) StoreNodeName(data map[string]interface{}) error {
	table := m.tables[utils.NodeName].(*tables.NodeName)
	monitorId := int64(data["monitorId"].(float64))
	sqlStr := table.GenerateCheckIdStr(monitorId)
	var count int
	m.db.QueryRow(sqlStr).Scan(&count)
	if count == 1 {
		return m.execUpdateNodeName(data) //update
	} else if count == 0 {
		insertSql := table.GenerateInsertStr(data)
		return m.execSql(insertSql)
	} else {
		deleteSql := table.GenerateDeleteStr(monitorId)
		m.execSql(deleteSql)
		insertSql := table.GenerateInsertStr(data)
		return m.execSql(insertSql)
	}
}

// 有节点则不更新，无节点则设置个空名称
func (m *Mysql) RegisterNode(monitorId int64) error {
	table := m.tables[utils.NodeName].(*tables.NodeName)
	sqlStr := table.GenerateCheckIdStr(monitorId)
	var count int
	data := make(map[string]interface{})
	data["monitorId"] = monitorId
	data["nodeName"] = ""
	m.db.QueryRow(sqlStr).Scan(&count)
	if count == 0 {
		insertSql := table.GenerateInsertStr(data)
		return m.execSql(insertSql)
	}
	return nil
}

func (m *Mysql) execUpdateNodeName(data map[string]interface{}) error {
	table := m.tables[utils.NodeName].(*tables.NodeName)
	sqlStr := table.GenerateUpdateStr(data)
	//fmt.Println("update node name:", sqlStr)
	res, err := m.db.Exec(sqlStr)
	if err != nil {
		return err
	} else {
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New("update data failed!")
			//failed
		} else {
			return nil
		}
	}
}
func (m *Mysql) showTables() []string {
	queryRes, _ := m.db.Query("SHOW TABLES")
	defer queryRes.Close()
	var result []string
	for queryRes.Next() {
		var table string
		queryRes.Scan(&table)
		result = append(result, table)
	}
	return result

}
func (m *Mysql) FetchOneNodeInfo(monitorId int64) (map[string]interface{}, error) {
	res, err := m.fetchNodeState(monitorId)
	res["nodeName"] = m.fetchNodeName(monitorId)
	return res, err
}
func (m *Mysql) fetchNodeState(monitorId int64) (map[string]interface{}, error) {
	sqlStr := fmt.Sprintf(`select * from %s where state->'$.monitorId'=%d order by state->'$.timestamp' desc limit 1;`,
		utils.GetColByAction(utils.NodeStates), monitorId)
	row := m.db.QueryRow(sqlStr)
	var id int64
	var data string
	err := row.Scan(&id, &data)
	/*
		err := row.Scan(&id, &data.MonitorId,
			&data.ViaGatewayId,
			&data.Rssi,
			&data.Snr,
			&data.Vbat,
			&data.AmbientTemperature,
			&data.ObjectTemperature,
			&data.Thermal,
			&data.Voltage,
			&data.Current,
			&data.Power,
			&data.Kwh,
			&data.Timestamp)*/
	res, _ := utils.JsonToMap(data)
	//res["lastUpdatedAt"] = res["timestamp"]
	//delete(res, "timestamp")
	if res["kwh"] != nil {
		res["status"] = "online"
	} else {
		res["status"] = "offline"
	}
	return res, err
}
func (m *Mysql) fetchNodeName(monitorId int64) string {
	//查询nodeName，如果没有，则为空
	sqlStr := fmt.Sprintf(`select nodeName from %s where monitorId=%d limit 1;`,
		utils.GetColByAction(utils.NodeName), monitorId)
	row := m.db.QueryRow(sqlStr)
	var nodeName string
	row.Scan(&nodeName)
	return nodeName
}

func (m *Mysql) FetchNodeEvents(action string, monitorId int64, count int64, page int64) ([]*map[string]interface{}, int64, error) {
	tableName := utils.GetColByAction(utils.Events)
	totalAmount := m.fetchEventsTotalAmountByMonitorId(action, monitorId)
	offSet := calculateOffset(totalAmount, count, page)
	actionFilter := ""
	if action == "" {
	} else {
		actionFilter = fmt.Sprintf(`and event->'$.action'="%s" `, action)

	}
	sqlStr := fmt.Sprintf(`select * from %s where event->'$.data.monitorId'=%d %sorder by event->'$.data.timestamp' desc limit %d offset %d;`,
		tableName, monitorId, actionFilter, count, offSet)
	var result []*map[string]interface{}
	rows, _ := m.db.Query(sqlStr)
	defer rows.Close()
	for rows.Next() {
		var id int64
		var dataStr string
		rows.Scan(&id, &dataStr)
		res, _ := utils.JsonToMap(dataStr)
		data := res["data"].(map[string]interface{})
		//timestamp to local
		timeStr := utils.TimeStampToLocalTimeStr(int64(data["timestamp"].(float64)))
		data["timestamp"] = timeStr
		res["data"] = data
		result = append(result, &res)
	}
	return result, totalAmount, nil
}

// select * from node_name s, node_states t where s.nodeName like '%test5%' and s.monitorId = t.state->'$.monitorId' and t.state->'$.thermal' is null;
func (m *Mysql) FetchBatchNodesInfo(count int64, page int64, name string, status string) ([]*map[string]interface{}, error) {
	//需要根据名称进行模糊查询，根据status online offline进行筛选
	totalAmount := m.fetchNodeCountByCondition(name, status)
	offset := calculateOffset(totalAmount, count, page)
	var result []*map[string]interface{}
	sqlStr := ""
	scanStates := func() {
		sqlStr = generateQueryBatchNodesStr(false, name, status)
		sqlStr = fmt.Sprintf(`%s order by state->'$.monitorId' asc limit %d offset %d;`, sqlStr, count, offset)
		rows, _ := m.db.Query(sqlStr)
		var id int64
		var data string
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&id, &data)
			res, _ := utils.JsonToMap(data)
			res["nodeName"] = m.fetchNodeName(int64(res["monitorId"].(float64)))
			if res["kwh"] != nil {
				res["status"] = "online"
			} else {
				res["status"] = "offline"
			}
			result = append(result, &res)
		}
	}
	scanNameAndStates := func() {
		sqlStr = generateQueryBatchNodesStr(false, name, status)
		sqlStr = fmt.Sprintf(`%s order by t.state->'$.monitorId' asc limit %d offset %d;`, sqlStr, count, offset)
		rows, _ := m.db.Query(sqlStr)
		var id int64
		var nodeName, data string
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&id, &id, &nodeName, &id, &data)
			res, _ := utils.JsonToMap(data)
			res["nodeName"] = nodeName
			if res["kwh"] != nil {
				res["status"] = "online"
			} else {
				res["status"] = "offline"
			}
			result = append(result, &res)
		}
	}
	if name != "" && status != "" {
		scanNameAndStates()
	} else if name == "" && status == "" {
		scanStates()
	} else {
		if name == "" {
			scanStates()
		} else {
			scanNameAndStates()
		}

	}
	/* 去重
	var finalRes []*map[string]interface{}
	for i := startId; i <= endId; i++ {
		monitorId := float64(i)
		count := 0
		for _, v := range result {
			if monitorId == (*v)["monitorId"].(float64) {
				count++
			}
		}
		if count == 1 {
			for _, v := range result {
				if monitorId == (*v)["monitorId"].(float64) {
					finalRes = append(finalRes, v)
				}
			}
		} else if count > 1 {
			//去重
			var dataById []*map[string]interface{}
			for _, v := range result {
				if monitorId == (*v)["monitorId"].(float64) {
					dataById = append(dataById, v)
				}
			}
			finalRes = append(finalRes, dataById[0])
		}
	}*/
	return result, nil
}

func calculateOffset(totalAmount int64, count int64, page int64) int64 {
	/*
		if count*page >= totalAmount {
			return 0
		} else {
			return count * page
		}*/
	return count * page
}
func (m *Mysql) fetchTotalAmountByMonitorId(tableName string, monitorId int64) int64 {
	var count int64 = 0
	sqlStr := fmt.Sprintf(`select count(*) from %s where monitorId=%d;`, tableName, monitorId)
	m.db.QueryRow(sqlStr).Scan(&count)
	return count
}
func (m *Mysql) fetchEventsTotalAmountByMonitorId(action string, monitorId int64) int64 {
	sqlStr := fmt.Sprintf(`select count(*) from %s where event->'$.data.monitorId'=%d`, utils.GetColByAction(utils.Events), monitorId)
	switch action {
	case utils.Dismount, utils.NFCCardRead, utils.MonitorSleeping:
		sqlStr = fmt.Sprintf(`%s and event->'$.action'="%s"`, sqlStr, action)
	default:
		sqlStr = sqlStr + ";"
	}
	var count int64 = 0
	m.db.QueryRow(sqlStr).Scan(&count)
	return count
}
func (m *Mysql) fetchTotalCount(tableName string) int64 {
	var count int64 = 0
	sqlStr := fmt.Sprintf(`select count(*) from %s;`, tableName)
	m.db.QueryRow(sqlStr).Scan(&count)
	return count
}
func (m *Mysql) fetchNodeCountByCondition(name string, status string) int64 {
	var count int64 = 0
	//select count(*) from node_states where json_extract(state,'$.thermal') is null;
	sqlStr := generateQueryBatchNodesStr(true, name, status)
	sqlStr = sqlStr + ";"
	m.db.QueryRow(sqlStr).Scan(&count)
	return count
}
func checkTableExist(list []string, table string) bool {
	for _, data := range list {
		if table == data {
			return true
		}
	}
	return false
}
