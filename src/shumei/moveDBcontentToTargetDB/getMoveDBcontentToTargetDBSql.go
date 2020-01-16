package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var sourceDB *sql.DB
var targetDB *sql.DB

var (
	class            *string
	variableIdPrefix *string
	sourceHost       *string
	sourceUser       *string
	sourcePasswd     *string
	sourcePort       *int
	sourceDbName     *string
	targetHost       *string
	targetUser       *string
	targetPasswd     *string
	targetPort       *int
	targetDbName     *string
)

func init() {
	setup()
}

func setup() {
	class = flag.String("class", "statistic-base", "counter define base-stats rule type")
	variableIdPrefix = flag.String("variableId", "common_YearMonthDay", "which rule be move")
	sourceUser = flag.String("sourceUser", "root", "source user")
	sourceHost = flag.String("sourceHost", "10.66.191.34", "source database host")
	sourcePasswd = flag.String("sourcePasswd", "shumeitest2018", "source database passwd")
	sourcePort = flag.Int("sourcePort", 3306, "mysql port")
	sourceDbName = flag.String("sourceDbName", "sentry", "source database")

	targetUser = flag.String("targetUser", "root", " target user")
	targetHost = flag.String("targetHost", "", "target database host")
	targetPasswd = flag.String("targetPasswd", "", "target database passwd")
	targetPort = flag.Int("targetPort", 3306, "mysql port")
	targetDbName = flag.String("targetDbName", "sentry", "target database")
	flag.Parse()
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
	connecttargetStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=1000ms&readTimeout=1000ms&writeTimeout=1000ms&charset=utf8",
		*targetUser,
		*targetPasswd,
		*targetHost,
		*targetPort,
		*targetDbName)
	if targetDB, err = sql.Open("mysql", connecttargetStr); err != nil {
		fmt.Println(fmt.Sprintf("Mysql Open error(%v)", err))
		return errors.New(fmt.Sprintf("Mysql Open error(%v)", err))
	}
	return nil
}

func scanAndMove() error {
	/*
		mysql> desc sentry_rule_engine_variable;
		+------------------+--------------+------+-----+-------------------+-----------------------------+
		| Field            | Type         | Null | Key | Default           | Extra                       |
		+------------------+--------------+------+-----+-------------------+-----------------------------+
		| id               | bigint(20)   | NO   | PRI | NULL              | auto_increment              |
		| variableId       | varchar(128) | NO   | MUL | NULL              |                             |
		| variableName     | varchar(256) | NO   |     | NULL              |                             |
		| definition       | text         | YES  |     | NULL              |                             |
		| class            | varchar(64)  | NO   |     | NULL              |                             |
		| type             | varchar(64)  | NO   |     | NULL              |                             |
		| organization     | varchar(64)  | NO   | MUL | NULL              |                             |
		| appId            | varchar(64)  | NO   |     | NULL              |                             |
		| enabled          | tinyint(4)   | YES  |     | 1                 |                             |
		| createTime       | bigint(20)   | NO   |     | 0                 |                             |
		| updateTime       | timestamp    | NO   |     | CURRENT_TIMESTAMP | on update CURRENT_TIMESTAMP |
		| eventId          | varchar(64)  | NO   |     | lending           |                             |
		| createUser       | varchar(64)  | NO   |     |                   |                             |
		| updateUser       | varchar(64)  | NO   |     |                   |                             |
		| statisticModel   | varchar(32)  | NO   |     |                   |                             |
		| definitionSource | text         | YES  |     | NULL              |                             |
		+------------------+--------------+------+-----+-------------------+-----------------------------+
		16 rows in set (0.00 sec)
	*/

	var (
		variableId       string
		variableName     string
		definition       string
		ruleType         string
		organization     string
		appId            string
		enabled          int
		eventId          string
		statisticModel   string
		definitionSource string
	)

	selectSql := "SELECT `variableId`, `variableName`, `definition`, `type`, `organization`,`appId`, `enabled`, `eventId`, `statisticModel`, `definitionSource` " +
		"FROM `sentry_rule_engine_variable` where `variableId` like '%" + *variableIdPrefix + "%';"

	//执行查询语句并获取结果
	rows, err := sourceDB.Query(selectSql)
	if err != nil {
		fmt.Println(fmt.Sprintf("Query failed.(%v)", err))
	}
	defer rows.Close()

	//逐条语句处理数据库内容
	count := 0
	for rows.Next() {
		//fmt.Print(fmt.Sprintf("moving current count = %v", count))

		rErr := rows.Scan(&variableId, &variableName, &definition, &ruleType, &organization, &appId, &enabled, &eventId, &statisticModel, &definitionSource)
		if rErr != nil {
			fmt.Println(fmt.Sprintf("Scan failed.(%v)", rErr))
			return rErr
		}

		var logSql = fmt.Sprintf("insert into `sentry_rule_engine_variable` (`variableId`, `variableName`, `definition`, `class`, `type`, `organization`,`appId`, `enabled`, `eventId`, `statisticModel`, `definitionSource`) values ('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v') on duplicate key update `enabled` = VALUES(enabled);",
			variableId,
			variableName,
			definition,
			*class,
			ruleType,
			organization,
			appId,
			enabled,
			eventId,
			statisticModel,
			definitionSource)
		fmt.Println(logSql)
		count++

	}
	return nil
}

func main() {
	fmt.Println("----- MOVE RULE TO ONLINE START. -----")

	InitDb()
	scanAndMove()

	fmt.Println("----- MOVE RULE TO ONLINE END. -----")

}
