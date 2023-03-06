package influx

import (
	"context"
	"fmt"
	"monitor-backend/utils"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var instance *Influx

type Influx struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
}

func InfluxInstance() *Influx {
	if instance == nil {
		instance = new(Influx)
		instance.init()
	}
	return instance
}
func (i *Influx) init() {

	i.client = influxdb2.NewClient(os.Getenv("INFLUX_URL"), os.Getenv("INFLUX_TOKEN"))
	i.writeAPI = i.client.WriteAPIBlocking(os.Getenv("INFLUX_ORG"), os.Getenv("INFLUX_BUCKET"))
	i.queryAPI = i.client.QueryAPI(os.Getenv("INFLUX_ORG"))
}

func (i *Influx) Insert(measurement string, tags map[string]string, data map[string]interface{}, ts time.Time) error {
	p := influxdb2.NewPoint(measurement, tags, data, ts)
	return i.writeAPI.WritePoint(context.Background(), p)
}
func (i *Influx) FindMonitorStateHistory(monitorId string, start int64, stop int64) ([]*map[string]interface{}, error) { //
	measurement := utils.GetColByAction(utils.MonitorState)
	//var start.stop int64 = 1671580790,1671580800

	queryStr := fmt.Sprintf(`from(bucket:"%s") |> range(start:%d, stop:%d) |> filter(fn: (r) => r._measurement == "%s" and r.monitorId == "%s") |> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")`,
		os.Getenv("INFLUX_BUCKET"), start, stop, measurement, monitorId)
	//fmt.Printf("influx queryStr:%s\n", queryStr)
	result, err := i.queryAPI.Query(context.Background(), queryStr)

	var list []*map[string]interface{}
	if err == nil {
		// Iterate over query response
		if result.Err() != nil {
			return list, result.Err()
		}
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				//fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			data := result.Record().Values()
			list = append(list, &data)
		}
		return list, err
	} else {
		return list, err
	}
}
func (i *Influx) FindLatestMonitorState(monitorId string) {

	//queryStr := `select * from monitor_state ORDER BY DESC LIMIT 1`
	//result, err := i.queryAPI.Query(context.Background(), queryStr)
}
