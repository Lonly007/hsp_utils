package hsp_utils

import (
	"fmt"
	"testing"
	// "time"
)

func _TestWriteJsonFiles(t *testing.T) {
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

	var x = RequestIteminfoST{}
	y := make([]OfflinePtInfoST, 1000)
	f1 := CreateRandFloat64List(0, 3, 1000)
	b1 := CreatRandBoolList(1000)
	for i := 0; i < 1000; i++ {
		y[i].SupportID = float64(i)
		y[i].Pos = f1[i]
		y[i].Speed = f1[i]
		y[i].Direction = b1[i]
		y[i].LeftHeight = f1[i]
		y[i].RightHeight = f1[i]
		y[i].PitchAngle = f1[i]
		y[i].Inclinator = f1[i]
	}
	x.ActType = "gaga"
	x.Datas = y

	var z interface{} = x

	err := WriteSTInfoJsonFiles(z, "D:\\release\\test.json")
	fmt.Println(err)
}

func _TestReadSTInfoJsonFiles(t *testing.T) {
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
	cfgInfo := "D:/release/test.json"

	x, err := ReadSTInfoJsonFiles_Example(cfgInfo)
	fmt.Println(x, err)
	y, err := ReadSTInfoJsonFiles(cfgInfo)
	fmt.Println(y, err)
}

func TestDelAndCheckFiles(t *testing.T) {
	cfgInfo := "D:/release/test.json"
	ex := CheckFileExist(cfgInfo)
	fmt.Println(ex)
	err := DelFiles(cfgInfo)
	fmt.Println(err)
}
