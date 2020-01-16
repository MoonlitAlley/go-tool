package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var (
	host      *string
	port      *string
	auth      *string
	bucket    *string
	redisconn redis.Conn
	err       error
)

func init() {
	host = flag.String("host", "", "redis host")
	port = flag.String("port", "", "redis port")
	auth = flag.String("auth", "", "redis password")
	bucket = flag.String("bucket", "", "redis hash key")
	flag.Parse()
}

func main() {
	//port = 6379
	redisconn, err = redis.DialTimeout("tcp", fmt.Sprintf("%s:%s", *host, *port),
		time.Duration(1000)*time.Millisecond,
		time.Duration(1000)*time.Millisecond,
		time.Duration(1000)*time.Millisecond)

	if err != nil {
		fmt.Print(fmt.Sprint("action=connect redis fail\tobj=redis\t"+"bucket="+"\terr= %s", err))
		return
	}
	if len(*auth) != 0 {
		_, err = redisconn.Do("AUTH", *auth)
		if err != nil {
			fmt.Print(fmt.Sprint("action=connect redis fail\tobj=redis\t"+"bucket="+"\terr= %s", err))
			return
		}
	}

	count, err := scanBucket(*bucket, 0, 0, redisconn)
	if err != nil {

	}
	result := fmt.Sprintf("redis key count =%d", count)
	fmt.Print(result)
}

func scanBucket(bucket string, itera int, count int, redisconn redis.Conn) (int, error) {
	defer func() (int, error) {
		if err := recover(); err != nil {
			return 0, errors.New(fmt.Sprintf("%v", err))
		}
		return count, nil
	}()

	ret, err := redisconn.Do("HSCAN", []interface{}{bucket, itera, "COUNT", 1000}...)
	if err != nil {
		return count, err
	}

	nextItera, _ := strconv.Atoi(string(ret.([]interface{})[0].([]byte)))

	kvList := (ret.([]interface{}))[1].([]interface{})

	count = count + len(kvList)/2

	if nextItera == 0 {
		return count, nil
	} else {
		countNext, errNext := scanBucket(bucket, nextItera, count, redisconn)
		if errNext != nil {
			return 0, nil
		} else {
			count = count + countNext
		}
	}
	return count, nil
}
