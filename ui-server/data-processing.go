package uiserver

import (
	"math"
	mysql "monitor-backend/mysql"
	"monitor-backend/utils"
	"os"
	"strconv"
	"time"
)

// Downsampling data by DOWN_SAMPLING 300s
func downsamplingData(list [][2]float64) [][2]float64 {

	downSampling, _ := strconv.ParseFloat(os.Getenv("DOWN_SAMPLING"), 64)
	if downSampling == 0 {
		downSampling = 300 //300s
	} else if downSampling == -1 {
		return list
	}
	res := [][2]float64{list[0]}
	startTimestamp := list[0][0]
	//startTime := time.Unix(int64(list[0][0]), 0)
	for _, data := range list {
		if data[0] > startTimestamp+downSampling {
			res = append(res, data)
			startTimestamp = startTimestamp + downSampling
		}
	}
	return res
}

// DUTY_EVERY 窗口300s
func getBodyDetectionData(list []*map[string]interface{}) [][2]float64 {
	every, _ := strconv.ParseInt(os.Getenv("DUTY_EVERY"), 10, 64)
	if every == 0 {
		every = 300
	}

	startTime := (*list[0])["_time"].(time.Time)
	var windowData []*map[string]interface{}
	var bodyDetection [][2]float64

	for _, data := range list {
		lastUpdatedAt := (*data)["_time"].(time.Time)
		windowTime := startTime.Add(time.Second * time.Duration(every))
		if !lastUpdatedAt.Before(windowTime) {
			//pointTime := startTime.Add(time.Second * time.Duration(every)).Unix()
			//pointValue := calculateDutyData(windowData)
			res := [2]float64{float64(windowTime.Unix()), calculateBodyDetectionData(windowData)}
			windowData = []*map[string]interface{}{} //clear data
			startTime = windowTime
			bodyDetection = append(bodyDetection, res)
		}
		windowData = append(windowData, data)
	}
	return downsamplingData(bodyDetection)
}

func calculateBodyDetectionData(list []*map[string]interface{}) float64 {
	rowCount, thermalCount, powerCount := float64(len(list)), 0.0, 0.0
	var objectTemperatureList []float64

	for _, data := range list {

		if (*data)["thermal"].(bool) {
			thermalCount++
		}
		if (*data)["current"].(float64) > 0.15 {
			powerCount++
		}
		objectTemperatureList = append(objectTemperatureList, (*data)["objectTemperature"].(float64))
		//fmt.Printf("t:%t c:%f obt:%f\n", (*data)["thermal"].(bool), (*data)["current"].(float64), (*data)["objectTemperature"].(float64))
	}
	variance := calculateVariance(objectTemperatureList)
	thermalPercent := thermalCount / rowCount
	powerPercent := powerCount / rowCount
	if thermalPercent < 0.1 {
		thermalPercent = thermalPercent * 10
	}
	variance = variance / 1000
	/*
		obT := ""
		for _, data := range objectTemperatureList {
			str := fmt.Sprintf("%f", data)
			obT += str + ","
		}
		fmt.Printf("总数:%d 方差:%f thermalPercent:%f powerPercent:%f 结果:%f 目标温度:%s\n", len(list), variance, thermalPercent, powerPercent, (variance+thermalPercent+powerPercent)/3, obT)
	*/
	if (variance+thermalPercent+powerPercent)/3 > 0.3 {
		return 1
	} else {
		return 0
	}
}
func calculateVariance(list []float64) float64 {
	var sum, s2 float64 = 0, 0
	for _, data := range list {
		sum += data
	}
	mean := sum / float64(len(list))
	for _, data := range list {
		s2 += math.Pow(data-mean, 2)
	}
	return s2 / float64(len(list))
}

/*
	func processOriginData(list []*map[string]interface{}) []*map[string]interface{} {
		res := []*map[string]interface{}{}
		for _, data := range list {
			(*data)["timestamp"] = (*data)["_time"].(time.Time).Unix()
			delete((*data), "result")
			res = append(res, data)
		}
		return res
	}
*/
func generateChartDataByKey(key string, list []*map[string]interface{}) [][2]float64 {
	var res [][2]float64
	for _, data := range list {
		timestamp := (*data)["_time"].(time.Time).Unix()
		var v = 0.0
		/*
			if key == "rssi" || key == "snr" {
				v = float64((*data)[key].(int64))
			} else {
				v = (*data)[key].(float64)
			}*/
		if key == "thermal" {
			if (*data)[key].(bool) {
				v = 1
			} else {
				v = 0
			}
		} else {
			v = (*data)[key].(float64)
		}

		res = append(res, [2]float64{float64(timestamp), v})
	}
	//return res
	return downsamplingData(res)
}
func getNodeEventLogs(monitorId int64, count int64, page int64, category int64) ([]*map[string]interface{}, int64) {
	res := []*map[string]interface{}{}
	var totalAmount int64
	switch category {
	case 0: //dismount
		res, totalAmount = fetchData(utils.Dismount, monitorId, count, page)
	case 1: //nfc
		res, totalAmount = fetchData(utils.NFCCardRead, monitorId, count, page)
	case 2: //sleeping
		res, totalAmount = fetchData(utils.MonitorSleeping, monitorId, count, page)
	case -1: //all
		res, totalAmount = fetchData("", monitorId, count, page)
	default:
	}
	return res, totalAmount
}
func getNodesInfo(count int64, page int64, name string, status string) ([]*map[string]interface{}, error) {
	return mysql.MysqlInstance().FetchBatchNodesInfo(count, page, name, status)
}
func fetchData(action string, monitorId int64, count int64, page int64) ([]*map[string]interface{}, int64) {
	events, totalAmount, _ := mysql.MysqlInstance().FetchNodeEvents(action, monitorId, count, page)
	return events, totalAmount
}
