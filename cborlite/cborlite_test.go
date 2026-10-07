package cborlite

import "testing"

func TestDecodeSmallInt(t *testing.T) {
	value, err := Decode([]byte{0x05}, &Counter{})
	if err != nil || value.(int) != 5 {
		t.Fatalf("应解出 5：%v %v", value, err)
	}
}

func TestDecodeText(t *testing.T) {
	value, _ := Decode([]byte{0x63, 'a', 'b', 'c'}, &Counter{})
	if value.(string) != "abc" {
		t.Fatalf("文本不对：%v", value)
	}
}

func TestDecodeArray(t *testing.T) {
	value, _ := Decode([]byte{0x82, 0x01, 0x02}, &Counter{})
	if len(value.([]any)) != 2 {
		t.Fatalf("数组长度不对：%v", value)
	}
}

func TestEncodeMap(t *testing.T) {
	data := EncodeMap(Map{Keys: []string{"a"}, Values: []any{1}}, &Counter{})
	if len(data) == 0 {
		t.Fatal("编码不该为空")
	}
}

func TestDecodeEof(t *testing.T) {
	if _, err := Decode(nil, &Counter{}); err == nil {
		t.Fatal("空输入应报错")
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	Decode([]byte{0x05}, c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
