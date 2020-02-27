package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var sourceDB *sql.DB
var fieldCountMap map[string]int64
var fieldCount uint64
var (
	sourceHost    *string
	sourceUser    *string
	sourcePasswd  *string
	sourcePort    *int
	sourceDbName  *string
	sourceDbTable *string
)

func init() {
	setup()
}

func setup() {
	sourceUser = flag.String("sourceUser", "smonline", "source user")
	sourceHost = flag.String("sourceHost", "10.66.191.34", "source database host")
	sourcePasswd = flag.String("sourcePasswd", "SMsmOnline2019", "source database passwd")
	sourcePort = flag.Int("sourcePort", 3306, "mysql port")
	sourceDbName = flag.String("sourceDbName", "profile_storage", "source database")
	sourceDbTable = flag.String("sourceDbTable", "storage_cluster", "source table")
	flag.Parse()
}

func main() {
	fmt.Println("----- GET FIELD COUNT START. -----")
	if dbErr := InitDb(); dbErr != nil {
		log.Println("InitDb Error", dbErr)
		return
	}
	defer sourceDB.Close()
	fieldCountMap = make(map[string]int64)
	if scanErr := scanAndCount(); scanErr != nil {
		log.Println("scan Error", scanErr)
		return
	}
	fmt.Println("----- GET FIELD COUNT END. -----")
	if len(fieldCountMap) > 0 {
		fmt.Println(fmt.Sprintf("keyCount = [%v], fieldCount =[%v]", fieldCount, len(fieldCountMap)))
		if countResult, marshalERR := json.Marshal(fieldCountMap); marshalERR != nil {
			fmt.Println("countResult marshal error")
		} else {
			fmt.Println(fmt.Sprintf("countResult = [%v]", string(countResult)))
		}
	} else {
		fmt.Println("NOT GET FIELD COUNT")
	}
	fmt.Println("----- finish. -----")
	return
}

func InitDb() error {
	connectSourceStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=1000ms&readTimeout=1000ms&writeTimeout=1000ms&charset=utf8",
		*sourceUser,
		*sourcePasswd,
		*sourceHost,
		*sourcePort,
		*sourceDbName)

	var err error
	if sourceDB, err = sql.Open("mysql", connectSourceStr); err != nil {
		fmt.Println(fmt.Sprintf("Mysql Open error(%v)", err))
		return errors.New(fmt.Sprintf("Mysql Open error(%v)", err))
	}
	return nil
}

func scanAndCount() error {
	/*
		mysql> desc storage_cluster;
		+-------------+--------------+------+-----+-------------------+-----------------------------+
		| Field       | Type         | Null | Key | Default           | Extra                       |
		+-------------+--------------+------+-----+-------------------+-----------------------------+
		| id          | bigint(20)   | NO   | PRI | NULL              | auto_increment              |
		| key         | varchar(128) | NO   | MUL | NULL              |                             |
		| originalKey | varchar(255) | NO   |     | NULL              |                             |
		| field       | varchar(128) | NO   |     | NULL              |                             |
		| value       | blob         | NO   |     | NULL              |                             |
		| expire      | bigint(20)   | NO   |     | NULL              |                             |
		| modifyTime  | timestamp    | NO   | MUL | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
		+-------------+--------------+------+-----+-------------------+-----------------------------+
		7 rows in set (0.00 sec)
	*/

	var (
		maxId       uint64
		originalKey string
		field       string
		value       string
		expire      string
		modifyTime  []uint8
	)

	//首先得到数据库当前表中的最大id
	getMaxIdSql := fmt.Sprintf("SELECT max(id) from `%s`;", *sourceDbTable)
	maxIdRows, err := sourceDB.Query(getMaxIdSql)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get Max Id Failed.(%v)", err))
	}
	maxIdRows.Next()
	maxIdSacnErr := maxIdRows.Scan(&maxId)
	if maxIdSacnErr != nil {
		fmt.Println(fmt.Sprintf("Scan Max Id Failed.(%v)", maxIdSacnErr))
	}
	maxIdRows.Close()

	//然后使用id区间对画像表进行扫描来统计数据
	for scanTimes := uint64(0); scanTimes*1000 < maxId; scanTimes++ {
		//输出当前执行进度
		load := ""
		for i := uint64(0); i < (scanTimes * 1000 * 100 / maxId); i++ {
			load = load + "="
			loadStr := fmt.Sprintf("[%s    %v]", load, i)
			fmt.Printf("\r%s", loadStr)
		}

		scanSql := fmt.Sprintf("SELECT `originalKey`, `field`, `value`, `expire`, `modifyTime` FROM `%s` where id >= %v and id < %v;",
			*sourceDbTable, scanTimes*1000, (scanTimes+1)*1000)
		dataRows, queryErr := sourceDB.Query(scanSql)
		if queryErr != nil {
			fmt.Println(fmt.Sprintf("Query failed.(%v); sql=(%v)", err, scanSql))
			continue
		}
		defer dataRows.Close()

		for dataRows.Next() {
			dataRowsSacnErr := dataRows.Scan(&originalKey, &field, &value, &expire, &modifyTime)
			if dataRowsSacnErr != nil {
				fmt.Println(fmt.Sprintf("DataRowsSacn Failed.(%v)", dataRowsSacnErr))
				continue
			}
			if _, ok := fieldCountMap[field]; ok {
				fieldCountMap[field] += 1
			} else {
				fieldCountMap[field] = 1
			}
			fieldCount++
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("")
	return nil
}
