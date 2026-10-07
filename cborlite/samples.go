package cborlite

import "fmt"

// Scenario 是一个固定场景。
type Scenario struct {
	Name string
	Run  func() map[string]any
}

func encodeText300() []byte {
	body := make([]byte, 300)
	out := []byte{0x79, 0x01, 0x2c}
	return append(out, body...)
}

// Samples 返回全部场景。
func Samples() []Scenario {
	return []Scenario{
		{Name: "len16", Run: func() map[string]any {
			value, _ := Decode(encodeText300(), &Counter{})
			text, ok := value.(string)
			if !ok {
				return map[string]any{"len": 0}
			}
			return map[string]any{"len": len(text)}
		}},
		{Name: "count200", Run: func() map[string]any {
			data := []byte{0x98, 0xc8}
			for i := 0; i < 200; i++ {
				data = append(data, 0x01)
			}
			value, _ := Decode(data, &Counter{})
			items, ok := value.([]any)
			if !ok {
				return map[string]any{"count": 0}
			}
			return map[string]any{"count": len(items)}
		}},
		{Name: "indef", Run: func() map[string]any {
			data := []byte{0x9f, 0x01, 0x02, 0x03, 0xff}
			value, _ := Decode(data, &Counter{})
			items, ok := value.([]any)
			if !ok {
				return map[string]any{"count": 0}
			}
			return map[string]any{"count": len(items)}
		}},
		{Name: "negint", Run: func() map[string]any {
			value, _ := Decode([]byte{0x20}, &Counter{})
			n, ok := value.(int)
			if !ok {
				return map[string]any{"value": 0}
			}
			return map[string]any{"value": n}
		}},
		{Name: "canon", Run: func() map[string]any {
			data := EncodeMap(Map{Keys: []string{"b", "a"}, Values: []any{1, 1}}, &Counter{})
			return map[string]any{"first": string(data[2:3])}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	m := &Map{}
	for i := 0; i < n; i++ {
		m.Keys = append(m.Keys, fmt.Sprintf("k%04d", i))
		m.Values = append(m.Values, i)
	}
	c := &Counter{}
	for i := 0; i < n; i++ {
		m.Get(fmt.Sprintf("k%04d", i), c)
	}
	return map[string]any{"elements": n, "scanned": c.Scanned}
}

var _ = fmt.Sprintf
