package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//array contain  character
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

//detect element contain and return index
func IsElementExist(s []string, str string) []int {
	positionArray := []int{}
	for i, v := range s {
		if v == str {
			positionArray = append(positionArray, i)
		}
	}
	return positionArray
}

//sum
func Sum(numbers ...int64) int64 {
	res := int64(0)
	for _, num := range numbers {
		res += num
	}
	return res
}

func JsonToMap(jsonStr string) (map[string]interface{}, error) {

	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}

	return m, nil
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// func GetWidth() uint {
// 	ws := &winsize{}
// 	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
// 		uintptr(syscall.Stdin),
// 		uintptr(syscall.TIOCGWINSZ),
// 		uintptr(unsafe.Pointer(ws)))
// 	if int(retCode) == -1 {
// 		panic(errno)
// 	}
// 	return uint(ws.Col)
// }

// func DivideLine(title string) error {
// 	hr := "-"
// 	width := GetWidth()
// 	divideLine := hr + strings.Repeat("-", int(width)-len(title)-3)
// 	divide := fmt.Sprintf(warningCode, title+divideLine)
// 	fmt.Println(divide)
// 	return nil
// }

func DivideLine(title string) error {

	divideLineCode := "\033[1;30;107m %s \033[0m"
	hr := "-"
	divideLine := hr + strings.Repeat("-", 100)
	divide := fmt.Sprintf(divideLineCode, title+divideLine)
	fmt.Println(divide)
	return nil
}

func MapToJson(param map[string][]string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

func CharacterRandomNumber(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func FormattimeToTimestampNanoStr(fromattime string) string {
	layout := "2006-01-0215:04:05"
	timepoint, err := time.ParseInLocation(layout, fromattime, time.Local)
	if err != nil {
		fmt.Println(err)
	}

	timestampNano := timepoint.UnixNano()
	timestampNanoStr := strconv.FormatInt(timestampNano, 10)
	return timestampNanoStr
}

func TimestampNanoStrToFormattime(timestamp string) string {
	unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
	t := time.Unix(0, unixIntValue)
	timeDisplay := t.Format("2006-01-02 15:04:05")
	return timeDisplay
}
