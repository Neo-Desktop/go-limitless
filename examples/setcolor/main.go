package main

import (
  "github.com/lucasb-eyer/go-colorful"
  "github.com/evq/go-limitless"
)

func main() {
  c := limitless.LimitlessController{}
  c.Host = "192.168.1.138"
  group := limitless.LimitlessGroup{}
  group.Id = 1
  group.Controller = &c
  c.Groups = []limitless.LimitlessGroup{group}

  color := colorful.Hsv(320.0, 1.0, 1.0)
  group.SendColor(color)
}
