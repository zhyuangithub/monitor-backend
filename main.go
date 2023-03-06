package main

import (
	datastorage "monitor-backend/data-storage"
	influx "monitor-backend/influxdb"
	mysql "monitor-backend/mysql"
	uiserver "monitor-backend/ui-server"

	"github.com/joho/godotenv"
	broadcast "github.com/teivah/broadcast"
)

func main() {

	godotenv.Load()
	//test()
	startServer()
}

func startServer() {
	relay := broadcast.NewRelay[string]()
	influx.InfluxInstance()
	mysql.MysqlInstance()
	go datastorage.DataStorageInstance().Init(relay)
	uiserver.UiServerInstance().Init(relay)
}
func test() {
	mysql.MysqlInstance()
	go datastorage.DataStorageInstance()
	uiserver.UiServerInstance()
}
