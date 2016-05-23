package main

import (
	"github.com/nesurion/go-limitless"
	"time"
)

func main() {
	c := limitless.LimitlessController{}
	c.Host = "192.168.2.141"
	group := limitless.LimitlessGroup{}
	group.Id = 1
	group.Controller = &c
	c.Groups = []limitless.LimitlessGroup{group}

	group.Disco()
	time.Sleep(1000 * time.Millisecond)
	group.DiscoSlower()
	time.Sleep(1000 * time.Millisecond)
	group.DiscoFaster()
}
