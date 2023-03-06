package mysql

type ITable interface {
	getTableName() string
	GenerateCreateTableStr() string
	GenerateInsertStr(data map[string]interface{}) string
}
