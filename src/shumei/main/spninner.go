package main

import (
	"fmt"
	"time"
)

func main() {

	load := ""
	for i := 0; i <= 100; i++ {
		load = load + "="
		loadStr := fmt.Sprintf("[%s    %v]", load, i)
		fmt.Printf("\r%s", loadStr)
		time.Sleep(1 * time.Second)
		if i == 4 {
			fmt.Println("xxx")
		}
	}
}
