package main

import (
	"github.com/nesurion/go-limitless"
)

func main() {
	c := limitless.LimitlessController{}
	c.Host = "192.168.2.141"
	group := limitless.LimitlessGroup{}
	group.Id = 1
	group.Controller = &c
	c.Groups = []limitless.LimitlessGroup{group}

	group.Night()
}
