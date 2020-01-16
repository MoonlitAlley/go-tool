// Copyright (c) 2019 SHUMEI Inc. All rights reserved.
// Authors: saifeiSong <songsaifei@ishumei.com>.

package main

import (
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"net"
	"os"
	"thrift/prediction"
	"time"
)

var host *string
var port *string
var requestId *string
var serviceId *string
var eventId *string
var organization *string
var appId *string
var tokenId *string
var data *string

func init() {
	setup()
}

func setup() {
	host = flag.String("host", "127.0.0.1", "host")
	port = flag.String("port", "6200", "port")
	requestId = flag.String("request_id", "987654321", "requestId")
	serviceId = flag.String("service_id", "POST_EVENT", "serviceId")
	eventId = flag.String("event_id", "register", "eventId")
	organization = flag.String("organization", "RlokQwRlVjUrTUlkIqOg", "organization")
	appId = flag.String("app_id", "default", "appId")
	tokenId = flag.String("token_id", "2rgde345v64y54y54t", "tokenId")
	data = flag.String("data", `{"eventId":"register","organization":"RlokQwRlVjUrTUlkIqOg","redis-counter":{"reg_token_distinct_phone_prefix_1d_detail":["RlokQwRlVjUrTUlkIqOg_hhh2015","RlokQwRlVjUrTUlkIqOg_hhh2014","RlokQwRlVjUrTUlkIqOg_hhh2013","RlokQwRlVjUrTUlkIqOg_hhh2012","RlokQwRlVjUrTUlkIqOg_hhh2010","RlokQwRlVjUrTUlkIqOg_hhh2000","RlokQwRlVjUrTUlkIqOg_hhh2007","RlokQwRlVjUrTUlkIqOg_hhh2004","RlokQwRlVjUrTUlkIqOg_hhh2008","RlokQwRlVjUrTUlkIqOg_hhh2011","RlokQwRlVjUrTUlkIqOg_hhh2003","RlokQwRlVjUrTUlkIqOg_hhh2005","RlokQwRlVjUrTUlkIqOg_hhh2002","RlokQwRlVjUrTUlkIqOg_hhh2006","RlokQwRlVjUrTUlkIqOg_hhh2001","RlokQwRlVjUrTUlkIqOg_hhh2009"]},"rule-engine":{"hits":[{"description":"高风险设备：账号异常聚集","descriptionV2":"高风险设备：账号异常聚集","model":"M006020115","riskLevel":"REJECT","score":800},{"description":"高风险手机号段：账号异常聚集","descriptionV2":"高风险手机号段：账号异常聚集","model":"M00604045","riskLevel":"REJECT","score":800},{"description":"高风险手机号段：账号异常聚集","descriptionV2":"高风险手机号段：账号异常聚集","model":"M00604046","riskLevel":"REJECT","score":800},{"description":"高风险设备：行为频度异常","descriptionV2":"高风险设备：行为频度异常","model":"M01104002","riskLevel":"REJECT","score":700},{"description":"高风险设备：账号异常聚集","descriptionV2":"高风险设备：账号异常聚集","model":"M02060101","riskLevel":"REJECT","score":900},{"description":"高风险设备：账号异常聚集","descriptionV2":"高风险设备：账号异常聚集","model":"M02060185","riskLevel":"REJECT","score":900},{"description":"高风险IP：账号异常聚集","descriptionV2":"高风险IP：账号异常聚集","model":"M02080101","riskLevel":"REJECT","score":800},{"description":"高风险手机号段：关联账号数","descriptionV2":"高风险手机号段：账号异常聚集","model":"M99210407","riskLevel":"REJECT","score":900},{"description":"同设备1d内关联的账号数量","descriptionV2":"高风险设备：账号异常聚集","model":"M99120102","riskLevel":"REJECT","score":800}],"model":"M02060101"}}`, "")
	flag.Parse()
}

func main() {
	timestamp := int64(time.Now().UnixNano() / 1000000)

	fmt.Println("----- PREDICTION START. -----")

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort(*host, *port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := prediction.NewPredictorClientFactory(useTransport, protocolFactory)

	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to "+*host+*port, " ", err)
		os.Exit(1)
	}

	defer transport.Close()

	request := prediction.NewPredictRequest()
	request.RequestId = requestId
	request.ServiceId = serviceId
	request.EventId = eventId
	request.AppId = appId
	request.Organization = organization
	request.TokenId = tokenId
	request.Timestamp = &timestamp
	request.Data = data

	result, err := client.Predict(request)

	// fmt.Println("result detail value:", result.GetDetail())

	fmt.Println(fmt.Sprintf("requestId=%v, result=%v", *requestId, *result))
}
