package surl

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	ErrLength   = errors.New("Invalid length")
	ErrKey      = errors.New("Invalid key")
	ErrId       = errors.New("Invalid id")
	ErrIdString = errors.New("Invalid id string")
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	LenSmall  = 4
	LenMiddle = 5
	LenLarge  = 6
)

type creator struct {
	length int
	hex    int
	key    string
	dic    []string
}

func NewCreator(length int) *creator {
	if length <= 0 || length >= 100 {
		panic("Invalid length")
	}
	return &creator{
		length: length,
		hex:    len(chars),
	}
}

func (c *creator) NewKey() (key string) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < c.length; i++ {
		bytes := []byte(chars)
		rand.Shuffle(len(bytes), func(i int, j int) { bytes[i], bytes[j] = bytes[j], bytes[i] })
		key += string(bytes)
	}
	return
}

func (c *creator) SetKey(key string) error {
	if len(key) != len(chars)*c.length {
		return ErrKey
	}
	c.key = key
	c.dic = c.dekey()
	return nil
}

func (c *creator) dekey() []string {
	strs := make([]string, 0)
	for i := 0; i < len(c.key); i += c.hex {
		strs = append(strs, c.key[i:i+c.hex])
	}
	return strs
}

func (c *creator) Encode(i int64) (string, error) {
	if i <= 0 {
		return "", ErrId
	}

	s := hexConvert(strconv.FormatInt(i, 10), 10, c.hex)

	for n, l := 0, c.length-len(s); n < l; n++ {
		s = "0" + s
	}

	if len(c.key) > 0 {
		s = c.upset(s)
	}

	return s, nil
}

func (c *creator) Decode(s string) (int64, error) {
	if len(s) > c.length {
		return 0, ErrIdString
	}

	if len(c.key) > 0 {
		s = c.reset(s)
	}

	s = strings.TrimLeft(s, "0")
	s = hexConvert(s, c.hex, 10)

	return strconv.ParseInt(s, 10, 0)
}

func (cr *creator) upset(s string) string {
	bytes := []byte(s)

	for i := cr.length - 2; i >= 0; i-- {
		bytes[i] = cr.byteShift(bytes[i], strings.IndexByte(cr.dic[i+1], bytes[i+1]))
	}
	bytes[cr.length-1] = cr.byteShift(bytes[cr.length-1], strings.IndexByte(cr.dic[0], bytes[0]))

	for k, v := range bytes {
		i := strings.IndexByte(chars, v)
		bytes[k] = cr.dic[k][i]
	}

	return string(bytes)
}

func (cr *creator) reset(s string) string {
	bytes := []byte(s)

	for k, v := range bytes {
		i := strings.IndexByte(cr.dic[k], v)
		bytes[k] = chars[i]
	}

	bytes[cr.length-1] = cr.byteShift(bytes[cr.length-1], -strings.IndexByte(cr.dic[0], bytes[0]))
	for i := 0; i <= cr.length-2; i++ {
		bytes[i] = cr.byteShift(bytes[i], -strings.IndexByte(cr.dic[i+1], bytes[i+1]))
	}

	return string(bytes)
}

func (cr *creator) byteShift(c byte, n int) byte {
	i := strings.Index(chars, string(c)) + n
	if n > 0 {
		for {
			if i < cr.hex {
				return chars[i]
			}
			i -= cr.hex
		}
	}
	if n < 0 {
		for {
			if i >= 0 {
				return chars[i]
			}
			i += cr.hex
		}
	}
	return c
}
