package uiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	influx "monitor-backend/influxdb"
	mysql "monitor-backend/mysql"
	"monitor-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// @Summary 获取数据记录
// @Schemes
// @Description 获取节点的数据记录
// @Tags Get
// @Param startTimestamp query int true "起始时间戳"
// @Param endTimestamp query int true "结束时间戳"
// @Param nodeId query int true "节点id"
// @Accept json
// @Produce json
// @Router /GetNodeDataHistory [get]
func getNodeDataHistoryHandler(c *gin.Context) {
	start, _ := c.GetQuery("startTimestamp")
	end, _ := c.GetQuery("endTimestamp")
	monitorIdStr, _ := c.GetQuery("nodeId")

	res := make(map[string]interface{})
	res["error"] = ""
	//res["timestamp"] = time.Now().Unix()
	res["time"] = time.Now().Format("2006-01-02 15:04:05")
	res["version"] = utils.Version
	res["every"] = os.Getenv("DUTY_EVERY") + "s"

	if start != "" && end != "" && monitorIdStr != "" {
		startTimestamp, _ := strconv.ParseInt(start, 10, 64)
		endTimestamp, _ := strconv.ParseInt(end, 10, 64)
		list, err := influx.InfluxInstance().FindMonitorStateHistory(monitorIdStr, startTimestamp, endTimestamp)
		if len(list) != 0 {
			res["bodyDetect"] = getBodyDetectionData(list)
			res["rssi"] = generateChartDataByKey("rssi", list)
			res["snr"] = generateChartDataByKey("snr", list)
			res["vbat"] = generateChartDataByKey("vbat", list)
			res["ambientTemperature"] = generateChartDataByKey("ambientTemperature", list)
			res["objectTemperature"] = generateChartDataByKey("objectTemperature", list)
			res["voltage"] = generateChartDataByKey("voltage", list)
			res["current"] = generateChartDataByKey("current", list)
			res["power"] = generateChartDataByKey("power", list)
			res["kwh"] = generateChartDataByKey("kwh", list)
			res["thermal"] = generateChartDataByKey("thermal", list)

			data, _ := json.Marshal(res)
			c.String(200, string(data))
		} else if err != nil {
			res["error"] = err.Error()
			data, _ := json.Marshal(res)
			c.String(200, string(data))
		} else {
			res["error"] = "No data found!"
			data, _ := json.Marshal(res)
			c.String(200, string(data))
		}
	} else {
		res["error"] = "Parameters not correct"
		data, _ := json.Marshal(res)
		c.String(200, string(data))
	}
}

// @Summary 获取节点信息
// @Schemes
// @Description 获取单个节点的信息，节点状态字符串种类：online,offline
// @Tags Get
// @Param nodeId query int true "节点id"
// @Accept json
// @Produce json
// @Router /GetNodeInfo [get]
func getNodeInfoHandler(c *gin.Context) {
	monitorIdStr, _ := c.GetQuery("nodeId")
	res := make(map[string]interface{})
	res["error"] = ""
	//res["timestamp"] = time.Now().Unix()
	//res["time"] = time.Now().Format("2006-01-02 15:04:05")
	res["version"] = utils.Version
	res["data"] = make(map[string]interface{})
	if monitorIdStr == "" {
		res["error"] = "No nodeId"
		resJson, _ := json.Marshal(res)
		c.String(200, string(resJson))
		return
	}
	monitorId, _ := strconv.ParseInt(monitorIdStr, 10, 64)
	data, err := mysql.MysqlInstance().FetchOneNodeInfo(monitorId)
	if err == nil {
		res["data"] = data
	} else {
		res["error"] = err.Error()
	}
	resJson, _ := json.Marshal(res)
	c.String(200, string(resJson))
}

// @Summary 批量获取节点信息
// @Schemes
// @Description 按页码和数量获取全部节点的信息，节点状态字符串种类：online,offline
// @Tags Get
// @Param count query int true "数量"
// @Param page query int true "页码"
// @Param name query string false "节点名称，可进行模糊查询"
// @Param status query string false "节点状态，online,offline两种"
// @Accept json
// @Produce json
// @Router /GetNodesInfo [get]
func nodesHandler(c *gin.Context) {
	countStr, _ := c.GetQuery("count")
	pageStr, _ := c.GetQuery("page")
	name, _ := c.GetQuery("name")
	status, _ := c.GetQuery("status")

	res := make(map[string]interface{})
	res["error"] = ""
	res["version"] = utils.Version
	res["data"] = []*map[string]interface{}{}
	if countStr == "" || pageStr == "" {
		res["error"] = "count or page is not correct!"
		resJson, _ := json.Marshal(res)
		c.String(200, string(resJson))
		return
	}

	count, _ := strconv.ParseInt(countStr, 10, 64)
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	if countStr == "" {
		count = 10
	}

	data, err := getNodesInfo(count, page, name, status)
	res["data"] = data
	if err != nil {
		res["error"] = err.Error()
	}
	resJson, _ := json.Marshal(res)
	c.String(200, string(resJson))
}

// @Summary 获取节点事件
// @Schemes
// @Description 获取单个节点的事件信息
// @Tags Get
// @Param nodeId query int true "节点id"
// @Param category query int false "事件类型分类 -1:all; 0:dismount; 1:nfc; 2:sleeping"
// @Param count query int true "数量"
// @Param page query int true "页码"
// @Accept json
// @Produce json
// @Router /GetNodeEventLogs [get]
func getNodeEventLogsHandler(c *gin.Context) {
	monitorIdStr, _ := c.GetQuery("nodeId")
	categoryStr, _ := c.GetQuery("category")
	countStr, _ := c.GetQuery("count")
	pageStr, _ := c.GetQuery("page")

	category, _ := strconv.ParseInt(categoryStr, 10, 64) //-1 all 0 dismount,1 nfc, 2 sleeping
	monitorId, _ := strconv.ParseInt(monitorIdStr, 10, 64)
	count, _ := strconv.ParseInt(countStr, 10, 64)
	page, _ := strconv.ParseInt(pageStr, 10, 64)

	res := make(map[string]interface{})
	res["error"] = ""
	res["version"] = utils.Version
	res["data"] = []*map[string]interface{}{}

	if monitorIdStr == "" {
		res["error"] = "No nodeId"
		resJson, _ := json.Marshal(res)
		c.String(200, string(resJson))
		return
	}
	if categoryStr == "" {
		category = -1
	}
	if countStr == "" {
		count = 10
	}
	res["data"], res["totalAmount"] = getNodeEventLogs(monitorId, count, page, category)
	resJson, _ := json.Marshal(res)
	c.String(200, string(resJson))
}

// @Summary websocket通知
// @Schemes
// @Description websocket通知接口
// @Tags Websocket
// @Accept json
// @Produce json
// @Router /notificationCenter [get]
func socketHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}, Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		fmt.Println("Error called:", reason)
	}} // use default options
	conn, err := upgrader.Upgrade(w, r, nil)
	conn.SetCloseHandler(func(code int, text string) error {
		return errors.New("close handler called")
	})
	if err != nil {
		fmt.Println("Error during connection upgradation:", err)
		return
	}
	runLoop := true
	closeCh := make(chan bool)
	go socketReadMessage(conn, &runLoop, closeCh)
	for runLoop {
		l := UiServerInstance().relay.Listener(1) // Create a listener with a buffer capacity of 1
		select {
		case res := <-l.Ch():
			//fmt.Println("ok:", ok)
			resMap, _ := utils.JsonToMap(res)
			resByte, _ := json.Marshal(resMap)
			err = conn.WriteMessage(1, resByte)
			//log.Printf("write message to:%s\n", conn.RemoteAddr())
			if err != nil {
				conn.Close()
				//fmt.Println("Error during message writing:", err)
				runLoop = false
				break
			}
		case <-closeCh:
			runLoop = false
			close(closeCh)
			break
		}
	}
	//fmt.Println("socketHandler end called")
}
func socketReadMessage(conn *websocket.Conn, runLoop *bool, ch chan bool) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			//fmt.Println("Error during message reading:", err)
			*runLoop = false
			ch <- true
			break
		}
	}
	//fmt.Println("socketReadMessage end called")
}

type changeNodeNameStruct struct {
	NodeName string `json:"nodeName"`
	NodeId   int64  `json:"nodeId"`
}

// @Summary 修改节点名称
// @Schemes
// @Description 修改某个节点的名称
// @Tags Post
// @Param Data body changeNodeNameStruct true "JSON数据"
// @Accept json
// @Produce json
// @Router /ChangeNodeName [post]
func changeNodeNameHandler(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	postData, err := utils.JsonToMap(string(body))
	res := make(map[string]interface{})
	res["error"] = ""
	if err != nil {
		res["error"] = "error with parsing data!"
		resJson, _ := json.Marshal(res)
		c.String(200, string(resJson))
		return
	}

	if postData["nodeId"] == nil || postData["nodeName"] == nil {
		res["error"] = "error with parsing nodeId or nodeName!"
		resJson, _ := json.Marshal(res)
		c.String(200, string(resJson))
		return
	}
	nodeName := postData["nodeName"].(string)
	nodeId := postData["nodeId"].(float64)
	if nodeName == "" || nodeId < 0 {
		res["error"] = "No node name or nodeId not correct!"
		resJson, _ := json.Marshal(res)
		c.String(200, string(resJson))
		return
	}
	//nodeId转为monitorId
	delete(postData, "nodeId")
	postData["monitorId"] = nodeId
	err = mysql.MysqlInstance().StoreNodeName(postData)
	resJson, _ := json.Marshal(res)
	c.String(200, string(resJson))
}
