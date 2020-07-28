package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	fd, _ :=  os.Open("~/Downloads/base_station_ip.txt")

	content, readErr := ioutil.ReadAll(fd)
	if readErr != nil {
		fmt.Println("read err", readErr)
	}

	ipList := make([]string, 0)
	if err:= json.Unmarshal(content, ipList); err != nil {
		fmt.Println("Unmarshal err", err)
	}


	for _, ip := range ipList {
		cmdStr := fmt.Sprintf("./bin/predict.sh  --port=7000 '{\"__request_type__\":\"profile\" ,\"data\":{\"ip\":\"%s\"}}'", ip)

		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		output, err := cmd.Output()
		if err != nil {
			log.Println(fmt.Printf("Execute Shell:%s failed with error:%s", cmdStr, err.Error()))
		}
		fmt.Println("get result", string(output))
	}

}



