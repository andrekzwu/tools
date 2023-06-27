package encode

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Struct to string
func Struct2String(param interface{}, pretty ...bool) string {
	if param == nil {
		return ""
	}
	body, _ := json.Marshal(param)
	if len(pretty) > 0 && pretty[0] {
		var buff bytes.Buffer
		_ = json.Indent(&buff, body, "", "\t")
		return buff.String()
	}
	return string(body)
}

// String2Struct
func String2Struct(body []byte, param interface{}) error {
	if err := json.Unmarshal(body, param); err != nil {
		return err
	}
	return nil
}

// String2uint32,conver string to uint32
func String2Uint32(str string) uint32 {
	value, _ := strconv.Atoi(str)
	return uint32(value)
}

// String2Int32
func String2Int32(str string) int32 {
	value, _ := strconv.Atoi(str)
	return int32(value)
}

// String2Int64
func String2Int64(str string) int64 {
	value, _ := strconv.ParseInt(str, 10, 64)
	return value
}

// Int642String
func Int642String(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Uint322String
func Uint322String(value uint32) string {
	return strconv.Itoa(int(value))
}
