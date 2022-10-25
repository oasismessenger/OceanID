package utils

import (
	"bytes"
	"encoding/binary"
)

func Int2Bytes(n int64) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, n)
	return buf.Bytes()
}

func Bytes2Int(b []byte) int64 {
	var n int64
	binary.Read(bytes.NewReader(b), binary.BigEndian, &n)
	return n
}
