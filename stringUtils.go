package hsp_utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case time.Time:
		const base_format = "2006-01-02 15:04:05"
		key = value.(time.Time).Format(base_format)
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

/**
 * 字符串首字母转化为大写 ios_bbbbbbbb -> Ios_bbbbbbbb
 */
func StrFirstToUpper(str string) string {
	var firstU string = str[:1]

	firstU = strings.ToUpper(firstU)

	return firstU + str[1:]
}

/**
 * 数据库字母转化为原格式 ios_bbbbbbbb -> IosBbbbbbbbb
 */
func StrDBToMapKey(str string) string {
	x := strings.Split(str, "_")
	resStr := ""
	for _, s := range x {
		resStr += StrFirstToUpper(s)
	}

	return resStr
}

func WriteSTInfoJsonFiles(stInfo interface{}, fileName string) error {
	// 创建文件
	filePtr, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return err
	}
	defer filePtr.Close()

	// 创建Json编码器
	// encoder := json.NewEncoder(filePtr)

	// err = encoder.Encode(srcDataList)
	// if err != nil {
	// 	fmt.Println("Encoder failed", err.Error())
	// 	return err
	// } else {
	// 	fmt.Println("Encoder success")
	// }

	// 带JSON缩进格式写文件
	data, err := json.MarshalIndent(stInfo, "", "  ")
	if err != nil {
		// fmt.Println("Encoder failed", err.Error())
		return err
	} else {
		// fmt.Println("Encoder success")
	}

	filePtr.Write(data)

	return nil
}

//示例程序结构体声明开始
type OfflinePtInfoST struct {
	SupportID   float64 `json:"supportID"`
	Pos         float64 `json:"pos"`
	Speed       float64 `json:"speed"`
	Direction   bool    `json:"direction"`
	LeftHeight  float64 `json:"leftHeight"`
	RightHeight float64 `json:"rightHeight"`
	PitchAngle  float64 `json:"pitchAngle"`
	Inclinator  float64 `json:"inclinator"`
}
type RequestIteminfoST struct {
	ActType     string            `json:"type"`
	Id          int               `json:"id"`
	StSID       int               `json:"stSID"`
	EndSID      int               `json:"endSID"`
	MaxSpd      int               `json:"maxSpd"`
	LeftHeight  int               `json:"leftHeight"`
	RightHeight int               `json:"rightHeight"`
	Datas       []OfflinePtInfoST `json:"datas"`
}

//示例程序结构体声明结束
func ReadSTInfoJsonFiles_Example(filePathName string) (x RequestIteminfoST, err error) {

	//json文件读取
	fd, err := os.Open(filePathName)
	if err != nil {
		fmt.Println("error when open file:", filePathName)
		return
	}
	defer fd.Close()
	err = json.NewDecoder(fd).Decode(&x)

	return
}

func ReadSTInfoJsonFiles(filePathName string) (x interface{}, err error) {

	//json文件读取
	fd, err := os.Open(filePathName)
	if err != nil {
		fmt.Println("error when open file:", filePathName)
		return
	}
	defer fd.Close()
	err = json.NewDecoder(fd).Decode(&x)

	return
}

func DelFiles(filePathName string) error {
	return os.Remove(filePathName)
}

func CheckFileExist(filePathName string) bool {
	var exist = true
	if _, err := os.Stat(filePathName); os.IsNotExist(err) {
		exist = false
	}

	return exist
}
