package limitless

import (
  "encoding/binary"
  "net"
  "bytes"
  "github.com/lucasb-eyer/go-colorful"
  "time"
)

type LimitlessController struct {
  Host string `json:"host"`
  Name string `json:"name"`
  Groups []LimitlessGroup `json:"groups"`
}

type LimitlessGroup struct {
  Id int `json:"id"`
  Type string `json:"type"`
  Name string `json:"name"`
  Controller *LimitlessController `json:"-"`
}

type LimitlessMessage struct {
  Key uint8
  Value uint8
  Suffix uint8
}

const (
  LIMITLESS_ADMIN_PORT = "48899"
  LIMITLESS_PORT = "8899"
)

const MAX_BRIGHTNESS = 0x1b

func NewLimitlessMessage() *LimitlessMessage {
  msg := LimitlessMessage{}
  msg.Suffix = 0x55
  return &msg
}

func (g *LimitlessGroup) SendColor(c colorful.Color) (error) {
  h, s, v := c.Hsv()
  h = 240.0 - h
  if h < 0 {
    h = 360.0 + h
  }
  scaled_h := uint8(h * 255.0 / 360.0)
  scaled_v := uint8(v * MAX_BRIGHTNESS)

  var err error

  if scaled_v < 0x02 {
    return g.Off()
  // If closer to white then a saturated color :D
  } else if s < 0.5 {
    err = g.White()
    if err != nil {
      return err
    }
  } else {
    err = g.Activate()
    if err != nil {
      return err
    }
    time.Sleep(100 * time.Millisecond)
    err = g.SetHue(scaled_h)
    if err != nil {
      return err
    }

  }
  err = g.Activate()
  if err != nil {
    return err
  }
  time.Sleep(100 * time.Millisecond)
  err = g.SetBri(scaled_v)
  return err
}

func (g *LimitlessGroup) SetHue(h uint8) (error) {
  msg := NewLimitlessMessage()
  msg.Key = 0x40
  msg.Value = h
  return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) SetBri(b uint8) (error) {
  //if b > MAX_BRIGHTNESS {
    //return err
  //}
  msg := NewLimitlessMessage()
  msg.Key = 0x4e
  msg.Value = b
  return g.Controller.sendMsg(msg)
}
func (g *LimitlessGroup) White() (error) {
  msg := NewLimitlessMessage()
  msg.Key = uint8(0xC5 + ((g.Id - 1) * 2))
  return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) On() (error) {
  msg := NewLimitlessMessage()
  msg.Key = uint8(0x45 + ((g.Id - 1) * 2))
  return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) Off() (error) {
  msg := NewLimitlessMessage()
  msg.Key = uint8(0x46 + ((g.Id - 1) * 2))
  return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) Activate() (error) {
  return g.On()
}

func (c *LimitlessController) sendMsg(msg *LimitlessMessage) (error) {
	conn, err := net.Dial("udp", c.Host + ":" + LIMITLESS_PORT)
  if err != nil {
    return err
  }
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msg)
  _, err = conn.Write(buf.Bytes())
  return err
}
