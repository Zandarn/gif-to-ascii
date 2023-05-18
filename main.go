package main

import (
	"time"
)

var Timeout = 50

func main() {
	g := New()
	go g.Print()

	time.Sleep(time.Duration(Timeout) * time.Second)
}
