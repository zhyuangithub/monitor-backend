package datastorage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"

	influx "monitor-backend/influxdb"
	mysql "monitor-backend/mysql"
	"monitor-backend/utils"

	"github.com/gin-gonic/gin"
	broadcast "github.com/teivah/broadcast"
)

var instance *DataStorage

type DataStorage struct {
	relay *broadcast.Relay[string]
}

func DataStorageInstance() *DataStorage {
	if instance == nil {
		instance = new(DataStorage)
		//instance.init()
	}
	return instance
}
func (d *DataStorage) Init(relay *broadcast.Relay[string]) {
	d.relay = relay
	d.startServer()
}

func (d *DataStorage) startServer() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		output := fmt.Sprintf("data storage server:%s %s\n", utils.Version, time.Now().Format("15:04:05"))
		c.String(200, output)
	})
	r.POST("", storageHandlers)
	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("DATASTORAGE_PORT"))
	r.Run(addr)
}

func storageHandlers(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	postData, err := utils.JsonToMap(string(body))
	if err != nil {
		c.String(200, "error with parsing data")
		return
	}
	action := postData["action"].(string)
	data := postData["data"].(map[string]interface{})
	//tags map[string]string
	tags := make(map[string]string)
	//mysql数据库内存的是东八区时间 带+00:00的自动转化为东八区时间
	//mysql会自动把字符串转为东八区timestamp变量
	timeStampStr := ""
	layoutStr := "2006-01-02T15:04:05"
	timeTsForWs := time.Now()
	if timeTsForWs.Location().String() == "UTC" {
		timeStampStr = timeTsForWs.Format(layoutStr) + "+00:00" //服务器的time.Now是UTC
	} else {
		timeStampStr = timeTsForWs.Format(layoutStr) + "+08:00"
	}

	var mySqlErr error
	switch action {
	case utils.GatewayState, utils.GatewayHeartBeat, utils.Dismount, utils.NFCCardRead:
		//dismount.nfccardread需转换为timestamp
		data["timestamp"] = timeStampStr
		if action == utils.Dismount || action == utils.NFCCardRead {
			data["timestamp"] = time.Now().Unix()
			res := make(map[string]interface{})
			res["notification"] = "onNodeEventLog"
			res["description"] = action
			wsData := make(map[string]interface{})
			if action == utils.Dismount {
				wsData["category"] = 0
			} else {
				wsData["category"] = 1
			}
			//0 dismount,1 nfc, 2 sleeping
			wsDetailData := cloneMaps(data)
			wsTimeStr := utils.TimeStampToLocalTimeStr(time.Now().Unix())
			//fmt.Printf("timeStr: %s\n", wsTimeStr)
			wsDetailData["timestamp"] = wsTimeStr
			wsData["data"] = wsDetailData
			res["data"] = wsData
			resByte, _ := json.Marshal(res)
			resStr := string(resByte)
			DataStorageInstance().relay.Broadcast(resStr)
		}
		err = mysql.MysqlInstance().InsertData(action, data)
	case utils.MonitorState:
		tags["monitorId"] = strconv.FormatFloat(math.Floor(data["monitorId"].(float64)), 'f', 0, 64)
		convertedTime, timeError := utils.TimeStrToTime(data["lastUpdatedAt"].(string))
		ts := time.Now()
		if timeError == nil {
			ts = convertedTime
		}
		mysqlData := cloneMaps(data)
		mysqlData["timestamp"] = ts.Unix()
		mySqlErr = mysql.MysqlInstance().StoreNodeStates(mysqlData)
		delete(data, "monitorId")
		err = influx.InfluxInstance().Insert(utils.GetColByAction(action), tags, data, ts)
	case utils.MonitorSleeping:
		//last updated是UTC+8,sleep保存的是json,时间转换为时间戳
		timeStamp, timeError := utils.TimeStrToTimestamp(data["lastUpdatedAt"].(string))
		if timeError != nil {
			timeStamp = time.Now().Unix()
		}
		data["timestamp"] = timeStamp
		err = mysql.MysqlInstance().InsertData(action, data)
		mysql.MysqlInstance().StoreNodeStates(data)
		res := make(map[string]interface{})

		res["notification"] = "onNodeEventLog"
		res["description"] = action
		wsData := make(map[string]interface{})
		wsData["category"] = 2
		wsDetailData := cloneMaps(data)
		//wsTimeStr := convertTimeStrToLocalStr(timeTsForWs)
		//fmt.Printf("timeStr: %s\n", wsTimeStr)
		delete(wsDetailData, "timestamp")
		wsData["data"] = wsDetailData
		res["data"] = wsData
		resByte, _ := json.Marshal(res)
		resStr := string(resByte)
		DataStorageInstance().relay.Broadcast(resStr)
	case utils.NodeName:
		err = mysql.MysqlInstance().StoreNodeName(data)
	default:
	}

	if err != nil {
		output := fmt.Sprintf("post action %s error:%s %s\n", action, err.Error(), timeStampStr)
		c.String(200, output)
	} else if mySqlErr != nil {
		output := fmt.Sprintf("post action %s error:%s %s\n", action, mySqlErr.Error(), timeStampStr)
		c.String(200, output)
	} else {
		output := fmt.Sprintf("post action %s succeed:%s\n", action, timeStampStr)
		c.String(200, output)
	}

}

func cloneMaps(tags map[string]interface{}) map[string]interface{} {
	cloneTags := make(map[string]interface{})
	for k, v := range tags {
		cloneTags[k] = v
	}
	return cloneTags
}
