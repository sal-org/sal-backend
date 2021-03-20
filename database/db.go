package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"

	LOGGER "salbackend/logger"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // for mysql driver
)

var db *sql.DB
var err error

// ConnectDatabase - connect to mysql database with given configuration
func ConnectDatabase() {
	db, err = sql.Open("mysql", CONFIG.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(CONFIG.DBConnectionPool)
	db.SetMaxIdleConns(CONFIG.DBConnectionPool)
	db.SetConnMaxLifetime(time.Hour)
}

// database utils

// InsertWithUniqueID - insert data into table with unique id
func InsertWithUniqueID(table string, uniqueDigits int, body map[string]string, id string) (string, string, bool) {
	var (
		status string
		ok     bool
	)
	for i := 0; i < CONSTANT.NumberOfTimesUniqueInserts; i++ {
		body[id] = generateRandomID(uniqueDigits)
		status, ok = InsertSQL(table, body)
		if !strings.EqualFold(status, CONSTANT.StatusCodeDuplicateEntry) {
			break
		}
	}
	return body[id], status, ok
}

func generateRandomID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CONSTANT.RandomIDDigits[rand.Intn(len(CONSTANT.RandomIDDigits))]
	}
	return string(b)
}

// RowCount - get number of items in database with specified query
func RowCount(tableName string, where string, args ...interface{}) int {
	data, _, _ := SelectProcess("select count(*) as ctn from "+tableName+" where "+where, args...)
	if len(data) == 0 {
		return 0
	}
	count, _ := strconv.Atoi(data[0]["ctn"])
	return count
}

// CheckIfExists - check if data exists in table
func CheckIfExists(table string, params map[string]string) bool {
	data, _, _ := SelectSQL(table, []string{"1"}, params)
	return len(data) > 0
}

// sql wrapper functions

// ExecuteSQL - execute statement with defined values
func ExecuteSQL(SQLQuery string, params ...interface{}) (sql.Result, error) {
	LOGGER.Log("ExecuteSQL", SQLQuery, params)
	return db.Exec(SQLQuery, params...)
}

// QueryRowSQL - get single data with defined values
func QueryRowSQL(SQLQuery string, params ...interface{}) string {
	LOGGER.Log("QueryRowSQL", SQLQuery, params)
	var value string
	db.QueryRow(SQLQuery, params...).Scan(&value)
	return value
}

// UpdateSQL - update data with defined values
func UpdateSQL(tableName string, params map[string]string, body map[string]string) (string, bool) {
	args := []interface{}{}

	if len(body) == 0 {
		return CONSTANT.StatusCodeBadRequest, false
	}
	SQLQuery := "update `" + tableName + "` set "

	init := false
	for key, val := range body {
		if init {
			SQLQuery += ","
		}
		SQLQuery += "`" + key + "` = ? "
		args = append(args, val)
		init = true
	}

	SQLQuery += " where "
	init = false
	for key, val := range params {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += "`" + key + "` = ? "
		args = append(args, val)
		init = true
	}

	LOGGER.Log("UpdateSQL", SQLQuery, args)

	_, err = db.Exec(SQLQuery, args...)
	if err != nil {
		fmt.Println("UpdateSQL", err)
		return CONSTANT.StatusCodeServerError, false // default
	}
	return CONSTANT.StatusCodeOk, true
}

// DeleteSQL - delete data with defined values
func DeleteSQL(tableName string, params ...map[string]string) (string, bool) {
	if len(params) == 0 {
		return CONSTANT.StatusCodeServerError, false // atleast one value should be specified for deleting, cannot delete all values
	}
	args := []interface{}{}

	SQLQuery := "delete from `" + tableName + "` where "

	init := false
	for key, val := range params[0] {
		if init {
			SQLQuery += " and "
		}
		SQLQuery += "`" + key + "` = ? "
		args = append(args, val[0])
		init = true
	}
	LOGGER.Log("DeleteSQL", SQLQuery, args)

	_, err = db.Exec(SQLQuery, args...)
	if err != nil {
		fmt.Println("DeleteSQL", err)
		return CONSTANT.StatusCodeServerError, false // default
	}
	return CONSTANT.StatusCodeOk, true
}

// InsertSQL - insert data with defined values
func InsertSQL(tableName string, body map[string]string) (string, bool) {
	if len(body) == 0 {
		return CONSTANT.StatusCodeBadRequest, false
	}
	SQLQuery, args := BuildInsertStatement(tableName, body)
	LOGGER.Log("InsertSQL", SQLQuery, args)

	_, err = db.Exec(SQLQuery, args...)
	if err != nil {
		fmt.Println("InsertSQL", err)
		return CONSTANT.StatusCodeServerError, false // default
	}
	return CONSTANT.StatusCodeCreated, true
}

// BuildInsertStatement - build insert statement with defined values
func BuildInsertStatement(tableName string, body map[string]string) (string, []interface{}) {
	args := []interface{}{}
	SQLQuery := "insert into `" + tableName + "` "
	keys := " ("
	values := " ("
	init := false
	for key, val := range body {
		if init {
			keys += ","
			values += ","
		}
		keys += " `" + key + "` "
		values += " ? "
		args = append(args, val)
		init = true
	}
	keys += ")"
	values += ")"
	SQLQuery += keys + " values " + values
	return SQLQuery, args
}

// SelectSQL - query data with defined values
func SelectSQL(tableName string, columns []string, params ...map[string]string) ([]map[string]string, string, bool) {
	args := []interface{}{}
	SQLQuery := "select " + strings.Join(columns, ",") + " from `" + tableName + "`"
	if len(params) > 0 {
		where := ""
		init := false
		for key, val := range params[0] {
			if init {
				where += " and "
			}
			where += " `" + key + "` = ? "
			args = append(args, val)
			init = true
		}
		if strings.Compare(where, "") != 0 {
			SQLQuery += " where " + where
		}
	}
	return SelectProcess(SQLQuery, args...)
}

// SelectProcess - execute raw select statement
func SelectProcess(SQLQuery string, params ...interface{}) ([]map[string]string, string, bool) {
	LOGGER.Log("SelectProcess", SQLQuery, params)

	rows, err := db.Query(SQLQuery, params...)
	if err != nil {
		fmt.Println("SelectProcess", err)
		return []map[string]string{}, CONSTANT.StatusCodeServerError, false // default
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("SelectProcess", err)
		return []map[string]string{}, CONSTANT.StatusCodeServerError, false // default
	}

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	data := []map[string]string{}
	rest := map[string]string{}
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		rest = map[string]string{}
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("SelectProcess", err)
			return []map[string]string{}, CONSTANT.StatusCodeServerError, false // default
		}

		for i, raw := range rawResult {
			if raw == nil {
				rest[cols[i]] = ""
			} else {
				rest[cols[i]] = string(raw)
			}
		}

		data = append(data, rest)
	}
	return data, CONSTANT.StatusCodeOk, true
}
