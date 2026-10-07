# cbor-lite

解码与规范编码 CBOR：附加信息四种宽度、不定长数组、负整数、映射键的字节序排序。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/cborlite --sample=len16
go run ./cmd/cborlite --sample=indef
go run ./cmd/cborlite --sample=negint
go run ./cmd/cborlite --sample=canon
go run ./cmd/cborlite --sample=work
```

## 口径

- **附加信息**：小于 24 是立即值；24/25/26/27 分别再读 1/2/4/8 字节大端；31 表示不定长。
- **不定长**：主类型 4 的不定长数组以 `0xFF` 结束。
- **负整数**：主类型 1 的值是 `-1 - n`。
- **规范序**：映射编码时键按各自编码后的字节序列升序排列。
- **一次遍历**：解码只从头到尾走一遍。

## 不变量

- 解出的整数在往返编码后不变。
- 不定长数组在 `0xFF` 处结束。
- `scanned` 不随元素数平方级放大：两千个元素的 `scanned` 不超过 6000。

## 输出契约

场景打印一行 JSON，键随场景而定。`--sample=work` 打印 `{"elements": N, "scanned": N}`。
