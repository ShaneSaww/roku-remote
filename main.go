package main

import (
	"fmt"
	"time"

	"github.com/ShaneSaww/roku-remote/roku"
)

var (
	searchTime = (5 * time.Second)
)

func main() {

	rokuIp, err := roku.GetRokuIp(searchTime)
	if err != nil {
		panic(err)
	}

	fmt.Println(rokuIp)
}
