package limitless

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/lucasb-eyer/go-colorful"
	"net"
	"time"
)

type LimitlessController struct {
	Host       string           `json:"host"`
	Name       string           `json:"name"`
	Connection net.Conn         `json:"-"`
	Groups     []LimitlessGroup `json:"groups"`
}

type LimitlessGroup struct {
	Id         int                  `json:"id"`
	Type       string               `json:"type"`
	Name       string               `json:"name"`
	Controller *LimitlessController `json:"-"`
}

type LimitlessMessage struct {
	Key    uint8
	Value  uint8
	Suffix uint8
}

const (
	LIMITLESS_ADMIN_PORT = "48899"
	LIMITLESS_PORT       = "8899"
)

const MAX_BRIGHTNESS = 0x1b

func NewLimitlessController(host string) (*LimitlessController, error) {
	c := LimitlessController{
		Host: host,
	}
	err := c.OpenConnection(host)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *LimitlessController) CloseConnection() error {
	err := c.Connection.Close()
	return err
}

func (c *LimitlessController) OpenConnection(host string) error {
	conn, err := net.Dial("udp", host+":"+LIMITLESS_PORT)
	if err == nil {
		c.Connection = conn
	}
	return err
}

func NewLimitlessMessage() *LimitlessMessage {
	msg := LimitlessMessage{}
	msg.Suffix = 0x55
	return &msg
}

func (m *LimitlessMessage) generateKey(hex int, g *LimitlessGroup) {
	m.Key = uint8(hex + ((g.Id - 1) * 2))
	return
}

func (g *LimitlessGroup) SendColorByte(code uint8) {
	msg := NewLimitlessMessage()
	msg.Key = 0x40
	msg.Value = code
	g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) SendColor(c colorful.Color) error {
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

func (g *LimitlessGroup) SetHue(h uint8) error {
	msg := NewLimitlessMessage()
	msg.Key = 0x40
	msg.Value = h
	err := g.Activate()
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) SetBri(b uint8) error {
	if b > MAX_BRIGHTNESS {
		return errors.New("brightness too high. (max 0x1B)")
	}
	msg := NewLimitlessMessage()
	msg.Key = 0x4E
	msg.Value = b
	err := g.Activate()
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) White() error {
	msg := NewLimitlessMessage()
	if g.Id == 0 {
		msg.Key = 0xC2
	} else {
		msg.generateKey(0xC5, g)
	}
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) On() error {
	msg := NewLimitlessMessage()
	if g.Id == 0 {
		msg.Key = 0x42
	} else {
		msg.generateKey(0x45, g)
	}
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) Off() error {
	msg := NewLimitlessMessage()
	if g.Id == 0 {
		msg.Key = 0x41
	} else {
		msg.generateKey(0x46, g)
	}
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) Night() error {
	msg := NewLimitlessMessage()
	if g.Id == 0 {
		msg.Key = 0xC1
	} else {
		msg.generateKey(0xC6, g)
	}
	err := g.Off()
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) Disco() error {
	msg := NewLimitlessMessage()
	msg.Key = 0x4D
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) DiscoFaster() error {
	msg := NewLimitlessMessage()
	msg.Key = 0x44
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) DiscoSlower() error {
	msg := NewLimitlessMessage()
	msg.Key = 0x43
	return g.Controller.sendMsg(msg)
}

func (g *LimitlessGroup) Activate() error {
	return g.On()
}

func (c *LimitlessController) sendMsg(msg *LimitlessMessage) error {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, msg)
	_, err := c.Connection.Write(buf.Bytes())
	return err
}
