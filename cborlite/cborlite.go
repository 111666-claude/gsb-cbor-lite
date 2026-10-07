// Package cborlite 解码与编码 CBOR。
package cborlite

import (
	"bytes"
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

	index map[string]int
}

func decode(data []byte, pos int, c *Counter) (any, int, error) {
	if pos >= len(data) {
		return nil, pos, errors.New("eof")
	}
	c.Scanned++
	head := data[pos]
	pos++
	major := head >> 5
	info := int(head & 0x1f)
	if info == 31 {
		if major != 4 {
			return nil, pos, errors.New("unsupported")
		}
		out := []any{}
		for {
			if pos >= len(data) {
				return nil, pos, errors.New("eof")
			}
			if data[pos] == 0xff {
				return out, pos + 1, nil
			}
			value, next, err := decode(data, pos, c)
			if err != nil {
				return nil, pos, err
			}
			out = append(out, value)
			pos = next
		}
	}
	arg, pos, err := readArg(data, pos, info)
	if err != nil {
		return nil, pos, err
	}
	switch major {
	case 0:
		return int(arg), pos, nil
	case 1:
		return -1 - int(arg), pos, nil
	case 2, 3:
		length := int(arg)
		if pos+length > len(data) {
			return nil, pos, errors.New("eof")
		}
		return string(data[pos : pos+length]), pos + length, nil
	case 4:
		out := make([]any, 0, arg)
		for i := 0; i < int(arg); i++ {
			value, next, err := decode(data, pos, c)
			if err != nil {
				return nil, pos, err
			}
			out = append(out, value)
			pos = next
		}
		return out, pos, nil
	}
	return nil, pos, errors.New("unsupported")
}

// readArg 按附加信息读出参数：24/25/26/27 再读 1/2/4/8 字节大端。
func readArg(data []byte, pos, info int) (uint64, int, error) {
	switch {
	case info < 24:
		return uint64(info), pos, nil
	case info == 24:
		if pos+1 > len(data) {
			return 0, pos, errors.New("eof")
		}
		return uint64(data[pos]), pos + 1, nil
	case info == 25:
		if pos+2 > len(data) {
			return 0, pos, errors.New("eof")
		}
		return uint64(binary.BigEndian.Uint16(data[pos:])), pos + 2, nil
	case info == 26:
		if pos+4 > len(data) {
			return 0, pos, errors.New("eof")
		}
		return uint64(binary.BigEndian.Uint32(data[pos:])), pos + 4, nil
	case info == 27:
		if pos+8 > len(data) {
			return 0, pos, errors.New("eof")
		}
		return binary.BigEndian.Uint64(data[pos:]), pos + 8, nil
	}
	return 0, pos, errors.New("unsupported")
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
	sort.Slice(pairs, func(i, j int) bool {
		return bytes.Compare(pairs[i].key, pairs[j].key) < 0
	})
	out := encodeUint(5, uint64(len(pairs)))
	for _, p := range pairs {
		c.Scanned++
		out = append(out, p.key...)
		out = append(out, p.val...)
	}
	return out
}

// Get 按键取映射里的值。
func (m *Map) Get(key string, c *Counter) (any, bool) {
	if m.index == nil {
		m.index = make(map[string]int, len(m.Keys))
		for i, k := range m.Keys {
			c.Scanned++
			m.index[k] = i
		}
	}
	c.Scanned++
	i, ok := m.index[key]
	if !ok {
		return nil, false
	}
	return m.Values[i], true
}

func encodeText(s string) []byte {
	out := encodeUint(3, uint64(len(s)))
	return append(out, []byte(s)...)
}

func encodeInt(n int) []byte {
	if n >= 0 {
		return encodeUint(0, uint64(n))
	}
	return encodeUint(1, uint64(-1-n))
}

func encodeUint(major byte, n uint64) []byte {
	head := major << 5
	switch {
	case n < 24:
		return []byte{head | byte(n)}
	case n <= 0xff:
		return []byte{head | 24, byte(n)}
	case n <= 0xffff:
		out := make([]byte, 3)
		out[0] = head | 25
		binary.BigEndian.PutUint16(out[1:], uint16(n))
		return out
	case n <= 0xffffffff:
		out := make([]byte, 5)
		out[0] = head | 26
		binary.BigEndian.PutUint32(out[1:], uint32(n))
		return out
	default:
		out := make([]byte, 9)
		out[0] = head | 27
		binary.BigEndian.PutUint64(out[1:], n)
		return out
	}
}
