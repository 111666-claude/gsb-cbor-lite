// Package cborlite 解码与编码 CBOR。
package cborlite

import (
	"encoding/binary"
	"errors"
	"sort"
)

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Map 是一个保持规范序的映射。
type Map struct {
	Keys   []string
	Values []any
}

func decode(data []byte, pos int, c *Counter) (any, int, error) {
	if pos >= len(data) {
		return nil, pos, errors.New("eof")
	}
	c.Scanned++
	head := data[pos]
	major := head >> 5
	info := int(head & 0x1f)
	pos++
	switch major {
	case 0:
		return int(info), pos, nil
	case 1:
		return -int(info), pos, nil
	case 2, 3:
		length := info
		if pos+length > len(data) {
			return nil, pos, errors.New("eof")
		}
		value := string(data[pos : pos+length])
		return value, pos + length, nil
	case 4:
		out := []any{}
		for i := 0; i < info; i++ {
			walk := pos
			for walk < pos {
				c.Scanned++
				walk++
			}
			value, next, err := decode(data, pos, c)
			if err != nil {
				return out, pos, err
			}
			out = append(out, value)
			pos = next
		}
		return out, pos, nil
	}
	return nil, pos, errors.New("unsupported")
}

// Decode 解码一个值。
func Decode(data []byte, c *Counter) (any, error) {
	value, _, err := decode(data, 0, c)
	return value, err
}

// EncodeMap 按规范序编码一个映射。
func EncodeMap(m Map, c *Counter) []byte {
	type pair struct {
		key []byte
		val []byte
	}
	pairs := make([]pair, len(m.Keys))
	for i, key := range m.Keys {
		pairs[i] = pair{key: encodeText(key), val: encodeInt(m.Values[i].(int))}
	}
	var out []byte
	out = append(out, 0xa0|byte(len(pairs)))
	for _, p := range pairs {
		c.Scanned++
		out = append(out, p.key...)
		out = append(out, p.val...)
	}
	return out
}

// Get 按键取映射里的值。
func (m *Map) Get(key string, c *Counter) (any, bool) {
	for i, k := range m.Keys {
		c.Scanned++
		if k == key {
			return m.Values[i], true
		}
	}
	return nil, false
}

func encodeText(s string) []byte {
	out := []byte{0x60 | byte(len(s))}
	return append(out, []byte(s)...)
}

func encodeInt(n int) []byte {
	if n >= 0 {
		return []byte{byte(n)}
	}
	return []byte{0x20 | byte(-n-1)}
}

var _ = sort.Ints
var _ = binary.BigEndian
