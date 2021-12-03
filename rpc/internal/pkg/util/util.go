package util

import (
	"encoding/json"
	"strconv"
)

// InterfaceToUint 失败会返回0
func InterfaceToUint(val interface{}) (value uint) {

	switch val.(type) {
	case int8:
		v, _ := val.(int8)
		value = uint(v)
		break
	case int32:
		v, _ := val.(int32)
		value = uint(v)
		break
	case int64:
		v, _ := val.(int64)
		value = uint(v)
		break
	case uint8:
		v, _ := val.(uint8)
		value = uint(v)
		break
	case uint32:
		v, _ := val.(uint32)
		value = uint(v)
		break
	case uint64:
		v, _ := val.(uint64)
		value = uint(v)
		break
	case float64:
		v, _ := val.(float64)
		value = uint(v)
		break
	case float32:
		v, _ := val.(float32)
		value = uint(v)
		break
	case string:
		v, _ := val.(string)
		va, _ := strconv.Atoi(v)
		value = uint(va)
		break
	case json.Number:
		v, _ := val.(json.Number)
		va, _ := strconv.Atoi(v.String())
		value = uint(va)
		break
	default:
		value = uint(0)
		break
	}

	return
}

// InterfaceToInt64 失败会返回0
func InterfaceToInt64(val interface{}) (value int64) {

	switch val.(type) {
	case int8:
		v, _ := val.(int8)
		value = int64(v)
		break
	case int32:
		v, _ := val.(int32)
		value = int64(v)
		break
	case int64:
		v, _ := val.(int64)
		value = int64(v)
		break
	case uint8:
		v, _ := val.(uint8)
		value = int64(v)
		break
	case uint32:
		v, _ := val.(uint32)
		value = int64(v)
		break
	case uint64:
		v, _ := val.(uint64)
		value = int64(v)
		break
	case float64:
		v, _ := val.(float64)
		value = int64(v)
		break
	case float32:
		v, _ := val.(float32)
		value = int64(v)
		break
	case string:
		v, _ := val.(string)
		va, _ := strconv.Atoi(v)
		value = int64(va)
		break
	case json.Number:
		v, _ := val.(json.Number)
		va, _ := strconv.Atoi(v.String())
		value = int64(va)
		break
	default:
		value = int64(0)
		break
	}

	return
}

// InterfaceToString 失败会返回""
func InterfaceToString(val interface{}) (value string) {

	switch val.(type) {
	case int:
		v, _ := val.(int)
		value = strconv.Itoa(v)
		break
	case int8:
		v, _ := val.(int8)
		value = strconv.Itoa(int(v))
		break
	case int32:
		v, _ := val.(int32)
		value = strconv.Itoa(int(v))
		break
	case int64:
		v, _ := val.(int64)
		value = strconv.Itoa(int(v))
		break
	case uint8:
		v, _ := val.(uint8)
		value = strconv.Itoa(int(v))
		break
	case uint32:
		v, _ := val.(uint32)
		value = strconv.Itoa(int(v))
		break
	case uint64:
		v, _ := val.(uint64)
		value = strconv.Itoa(int(v))
		break
	case float64:
		v, _ := val.(float64)
		value = strconv.Itoa(int(v))
		break
	case float32:
		v, _ := val.(float32)
		value = strconv.Itoa(int(v))
		break
	case string:
		v, _ := val.(string)
		value = v
		break
	case json.Number:
		v, _ := val.(json.Number)
		value = string(v)
		break
	default:
		value = string("")
		break
	}

	return
}
