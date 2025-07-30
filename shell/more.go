package shell

import (
	"fmt"
	"lokishell/kernel"
	"lokishell/util"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/eiannone/keyboard"
)

func checkIpExisted() bool {
	for key, valueList := range filterMap {
		if strings.Contains(key, "ip") && len(valueList) > 0 {
			return true
		}
	}
	return false
}

func checkTimeRangeBoundary(starttime string, endtime string, day int) bool {
	timepointTemp, _ := time.ParseInLocation(layout, starttime, time.Local)
	endtimeOneday := timepointTemp.Add(+time.Hour * time.Duration(24*day)).Format(layout)
	return endtimeOneday >= endtime
}

func correlationCheck(starttime string, endtime string, keyword string, isGrepContext bool) bool {
	// 判断有grep参数，但没有keyword
	if isGrepContext && (keyword == "nokey" || keyword == "") {
		return false
	}
	// 判断时间范围是否超过1天
	if !checkTimeRangeBoundary(starttime, endtime, 1) {
		// fmt.Println("目前暂时不支持查询日期范围超过7天,请缩小日期查询范围")
		fmt.Println("目前暂时不支持该语法查询,或者时间范围超过24小时")
		return false
	}
	// 判断 app_name 下的日志量级别：low,mid,high
	log_scale := kernel.FetchLogScaleByApp(app_name, token)
	switch log_scale {
	case "mid":
		if !checkIpExisted() {
			fmt.Println("该应用日志量级别为(中),目前暂时无法直接查询所有服务,请先使用[filter]命令添加[_ip_],再进行查询")
			filterdoc()
			return false
		}
	case "high":
		if !checkIpExisted() {
			fmt.Println("该应用日志量级别为(高),目前暂时不支持直接查询所有服务,请先使用[filter]命令添加[_ip_],再进行查询")
			filterdoc()
			return false
		}
		if !checkTimeRangeBoundary(starttime, endtime, 1) {
			fmt.Println("该应用日志量级别为(高),目前暂时不支持查询日期范围超过1天,请缩小查询日期范围为1天内")
			return false
		}
	}
	return true
}

func more(l *readline.Instance, filterMap map[string][]string, starttime string, endtime string, step string, limit string, keyword string, isPositiveSequence bool, timepoint bool, isGrepContext bool, isRelativeTime bool) bool {
	isIntegralPoint := false
	if endtime == "jt" || endtime == "zt" || endtime == "qt" || endtime == "bz" || endtime == "sz" || endtime == "by" || endtime == "sy" ||
		strings.HasSuffix(endtime, "fz") || strings.HasSuffix(endtime, "xs") {
		isIntegralPoint = true
		timepoint = false
		isRelativeTime = false
	}

	if !isRelativeTime && !isIntegralPoint {
		if !correlationCheck(starttime, endtime, keyword, isGrepContext) {
			return false
		}
	}
	lmt, err := strconv.Atoi(limit)
	if err != nil {
		return false
	}
	direction := "FORWARD"
	if !isPositiveSequence {
		direction = "BACKWARD"
	}

	var displayLayout = "2006-01-02 15:04:05"
	//根据时间点查找日志
	if timepoint {

		timeData := 60
		var err error
		if starttime != endtime {
			minStr := endtime[:len(endtime)-3]
			timeData, err = strconv.Atoi(minStr)
			if err != nil {
				timeData = 60
			}
		}
		timepoint, _ := time.ParseInLocation(layout, starttime, time.Local)
		starttime1 := timepoint.Add(-time.Minute * time.Duration(timeData)).Format(layout)
		endtime1 := timepoint.Format(layout)
		starttime2 := timepoint.Format(layout)
		endtime2 := timepoint.Add(+time.Minute * time.Duration(timeData)).Format(layout)

		fetchEvent(filterMap, starttime1, endtime1, step, limit, keyword, direction, true, isGrepContext)
		util.DivideLine(timepoint.Format(displayLayout))
		fetchEvent(filterMap, starttime2, endtime2, step, limit, keyword, direction, false, isGrepContext)
		return true
	}

	startTimestampStr := ""
	endTimestampStr := ""
	starttimeDisplay := starttime
	endtimeDisplay := endtime
	switch {
	case isIntegralPoint:
		switch {
		case strings.HasSuffix(endtime, "fz"):
			fzStr := endtime[:len(endtime)-2]
			fzunit, _ := strconv.Atoi(fzStr)

			startTimestampStr, endTimestampStr, starttimeDisplay, endtimeDisplay = calculateTimeRange(fzunit, displayLayout, "fz")
		case strings.HasSuffix(endtime, "xs"):
			xsStr := endtime[:len(endtime)-2]
			xsunit, _ := strconv.Atoi(xsStr)

			startTimestampStr, endTimestampStr, starttimeDisplay, endtimeDisplay = calculateTimeRange(xsunit, displayLayout, "xs")
		case endtime == "jt":
			startTimestampStr, endTimestampStr, starttimeDisplay, endtimeDisplay = calculateTimeRange(0, displayLayout, "jt")
		case endtime == "zt":
			startTimestampStr, endTimestampStr, starttimeDisplay, endtimeDisplay = calculateTimeRange(0, displayLayout, "zt")
		case endtime == "qt":
			startTimestampStr, endTimestampStr, starttimeDisplay, endtimeDisplay = calculateTimeRange(0, displayLayout, "qt")
		}
	case isRelativeTime:
		timeunit := endtime

		currenttime := time.Now()
		currentTimestamp := currenttime.UnixNano()
		endTimestampStr = strconv.FormatInt(currentTimestamp, 10)
		endtimeDisplay = currenttime.Format(displayLayout)
		var starttime time.Time
		if strings.HasSuffix(endtime, "m") || strings.HasSuffix(endtime, "h") {
			timesub, _ := time.ParseDuration("-" + timeunit)
			starttime = currenttime.Add(timesub)
		}
		if strings.HasSuffix(endtime, "h") {
			hunit := endtime[:len(endtime)-1]
			hint, _ := strconv.Atoi(hunit)
			if hint > 24 {
				fmt.Println("抱歉,目前系统暂时不支持查询超过24小时")
				return false
			}
		}
		if strings.HasSuffix(endtime, "d") {
			dStr := endtime[:len(endtime)-1]
			dunit, _ := strconv.Atoi(dStr)
			starttime = currenttime.Add(-time.Hour * time.Duration(dunit*24))
		}
		startTimestamp := starttime.UnixNano()
		startTimestampStr = strconv.FormatInt(startTimestamp, 10)
		starttimeDisplay = starttime.Format(displayLayout)
	default:
		startTimestampStr = util.FormattimeToTimestampNanoStr(starttime)
		endTimestampStr = util.FormattimeToTimestampNanoStr(endtime)
	}

	isMorePaging := false
	forwardrownum := fetchEvent(filterMap, startTimestampStr, endTimestampStr, step, limit, keyword, direction, true, isMorePaging)
	//aliyun集群和upyun集群响应行数都小于limit时，自动退出查询
	if forwardrownum["aliyuntotalrownum"].(int) < lmt && forwardrownum["upyuntotalrownum"].(int) < lmt && forwardrownum["totalnum"].(int) < lmt {
		// fmt.Println("该时间段内已无新数据")
		return true
	}
	promt_msg := "--More-- Press Space to Page Down/Press ESC or CtrlC to quit"

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	promt := fmt.Sprintf(pageCode, promt_msg)
	timepromt := fmt.Sprintf(pageCode, starttimeDisplay+"~"+endtimeDisplay)

	fmt.Println(promt + "\t" + timepromt)
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		//fmt.Printf("You pressed: rune %q, key %X\r\n", char, key)
		if key == keyboard.KeySpace {

			startTimePageunitStr := ""
			endTimePageunitStr := ""
			switch direction {
			case "FORWARD":
				startTimePageunitStr = forwardrownum["maxLastTimestamp"].(string)
			case "BACKWARD":
				endTimePageunitStr = forwardrownum["minFirstTimestamp"].(string)
			}

			// startTimestampInt, _ := strconv.ParseInt(startTimePageunitStr, 10, 64)
			// startTimePage := time.Unix(0, startTimestampInt)
			// endtimePage := startTimePage.Add(+time.Minute * 10)

			// endtime, _ := time.ParseInLocation(layout, endtime, time.Local)
			// endtimePageunixInt := endtimePage.UnixNano()
			// endtimeunixInt := endtime.UnixNano()
			// endtimeunixStr := strconv.FormatInt(endtimeunixInt, 10)

			// starttimePageFormat := startTimePage.Format(displayLayout)
			// endtimePageFormat := endtimePage.Format(displayLayout)
			// fmt.Println(starttimePageFormat)
			// fmt.Println(endtimePageFormat)

			//not allow to extend endtime
			// if endtimePageunixInt > endtimeunixInt {
			// 	endtimePageunixInt = endtimeunixInt
			// }
			if startTimePageunitStr >= endTimestampStr || (direction == "BACKWARD" && endTimePageunitStr <= startTimestampStr) {
				promt_msg = "--END--"
				promt = fmt.Sprintf(pageCode, promt_msg)
				fmt.Println(promt)
				break
			} else {
				isMorePaging = true
				switch direction {
				case "FORWARD":
					forwardrownum = fetchEvent(filterMap, startTimePageunitStr, endTimestampStr, step, limit, keyword, direction, true, isMorePaging)
				case "BACKWARD":
					forwardrownum = fetchEvent(filterMap, startTimestampStr, endTimePageunitStr, step, limit, keyword, direction, true, isMorePaging)
				}
				if forwardrownum["aliyuntotalrownum"].(int) < lmt && forwardrownum["upyuntotalrownum"].(int) < lmt && forwardrownum["totalnum"].(int) < lmt {
					// fmt.Println("该时间段内已无新数据")
					return true
				}

				promt := fmt.Sprintf(pageCode, promt_msg)
				fmt.Println(promt + "\t" + timepromt)
			}

		}
		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			break
		}
	}

	return true

}

func calculateTimeRange(timevalue int, displayLayout string, pointtype string) (startTimestampStr string, endTimestampStr string, starttimeDisplay string, endtimeDisplay string) {
	var startTime, endTime time.Time
	dateNow := time.Now()
	switch {
	case pointtype == "fz":
		endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), dateNow.Hour(), dateNow.Minute(), 0, 0, dateNow.Location())
		startTime = endTime.Add(-time.Minute * time.Duration(timevalue))
	case pointtype == "xs":
		endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), dateNow.Hour(), 0, 0, 0, dateNow.Location())
		startTime = endTime.Add(-time.Hour * time.Duration(timevalue))
	default:
		switch pointtype {
		case "zt":
			dateNow = dateNow.AddDate(0, 0, -1)
		case "qt":
			dateNow = dateNow.AddDate(0, 0, -2)
		}
		startTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, dateNow.Location())
		endTime = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 23, 59, 59, 0, dateNow.Location())
	}

	startTimestamp := startTime.UnixNano()
	endTimestamp := endTime.UnixNano()
	startTimestampStr = strconv.FormatInt(startTimestamp, 10)
	endTimestampStr = strconv.FormatInt(endTimestamp, 10)
	starttimeDisplay = startTime.Format(displayLayout)
	endtimeDisplay = endTime.Format(displayLayout)
	return
}
