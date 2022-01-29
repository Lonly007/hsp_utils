package hsp_utils

import (
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type HspRWLock struct {
	sync.RWMutex
}

func CreateRandFloat64List(min, max, cnt int) (res []float64) {
	if cnt <= 0 {
		return nil
	} else {
		res = make([]float64, cnt)
	}

	rand.Seed(time.Now().UnixNano()) // 纳秒时间戳

	offset := int(math.Abs(float64(min-max))) + 1
	for i := 0; i < cnt; i++ {
		res[i] = (float64)(rand.Intn(offset)) + rand.Float64() + math.Min((float64)(min), (float64)(max))
	}

	return
}

func CreateRandIntList(min, max, cnt int) (res []int) {
	if cnt <= 0 {
		return nil
	} else {
		res = make([]int, cnt)
	}

	rand.Seed(time.Now().UnixNano()) // 纳秒时间戳

	offset := int(math.Abs(float64(min-max))) + 1
	for i := 0; i < cnt; i++ {
		res[i] = rand.Intn(offset) + int(math.Min((float64)(min), (float64)(max)))
	}

	return
}

func CreatRandBoolList(cnt int) (res []bool) {
	if cnt <= 0 {
		return nil
	} else {
		res = make([]bool, cnt)
	}

	rand.Seed(time.Now().UnixNano()) // 纳秒时间戳

	tmpX := 0
	for i := 0; i < cnt; i++ {
		tmpX = rand.Intn(100)
		res[i] = (tmpX >= 50)
	}

	return
}

func ConvInterface2Float(srcData interface{}) float64 {
	var res float64
	switch srcData.(type) {
	case float64:
		res = srcData.(float64)
	case float32:
		res = float64(srcData.(float32))
	case int:
		res = float64(srcData.(int))
	case uint:
		res = float64(srcData.(uint))
	case int8:
		res = float64(srcData.(int8))
	case uint8:
		res = float64(srcData.(uint8))
	case int16:
		res = float64(srcData.(int16))
	case uint16:
		res = float64(srcData.(uint16))
	case int32:
		res = float64(srcData.(int32))
	case uint32:
		res = float64(srcData.(uint32))
	case int64:
		res = float64(srcData.(int64))
	case uint64:
		res = float64(srcData.(uint64))
	case string:
		res, _ = strconv.ParseFloat(srcData.(string), 64)
	default:
		res = 0.0
	}

	return res
}

func ConvInterface2Int64(srcData interface{}) int64 {
	var res int64
	switch srcData.(type) {
	case float64:
		res = int64(srcData.(float64))
	case float32:
		res = int64(srcData.(float32))
	case int:
		res = int64(srcData.(int))
	case uint:
		res = int64(srcData.(uint))
	case int8:
		res = int64(srcData.(int8))
	case uint8:
		res = int64(srcData.(uint8))
	case int16:
		res = int64(srcData.(int16))
	case uint16:
		res = int64(srcData.(uint16))
	case int32:
		res = int64(srcData.(int32))
	case uint32:
		res = int64(srcData.(uint32))
	case int64:
		res = int64(srcData.(int64))
	case uint64:
		res = int64(srcData.(uint64))
	case string:
		res, _ = strconv.ParseInt(srcData.(string), 10, 64)
	default:
		res = 0
	}

	return res
}

func ConvInterface2Bool(srcData interface{}) bool {
	var res bool
	switch srcData.(type) {
	case float64:
		res = true
	case float32:
		res = true
	case int:
		res = ((srcData.(int)) != 0)
	case uint:
		res = (srcData.(uint) != 0)
	case int8:
		res = (srcData.(int8) != 0)
	case uint8:
		res = (srcData.(uint8) != 0)
	case int16:
		res = (srcData.(int16) != 0)
	case uint16:
		res = (srcData.(uint16) != 0)
	case int32:
		res = (srcData.(int32) != 0)
	case uint32:
		res = (srcData.(uint32) != 0)
	case int64:
		res = (srcData.(int64) != 0)
	case uint64:
		res = (srcData.(uint64) != 0)
	case bool:
		res = srcData.(bool)
	case string:
		res = (srcData.(string) == "true")
	default:
		res = false
	}

	return res
}

func GetCurrentDir() (string, error) {
	ex, err := os.Executable()

	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	realPath, err := filepath.EvalSymlinks(exPath)
	if err != nil {
		return "", err
	}

	return realPath, nil
}
