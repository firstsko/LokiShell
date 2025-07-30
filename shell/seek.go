package shell

import (
	"fmt"
	"lokishell/kernel"
	"lokishell/util"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

var lastrownum = 5
var lasttaginfo = ""
var maxLastTimestamp = ""
var threadglobal = ""

func seek(id string, mode string, traceid string, border string, upnum string, downnum string) bool {
	if id == "" {
		help(false)
		return false
	}
	// uuid, err1 := uuid.FromString(id)
	// if err1 != nil {
	// 	fmt.Println("failed to parse UUID")
	// 	help(false)
	// 	return false
	// }
	upnumint, err2 := strconv.Atoi(upnum)
	if err2 != nil {
		fmt.Println("failed to parse the above lines")
		help(false)
		return false
	}
	downnumint, err3 := strconv.Atoi(downnum)
	if err3 != nil {
		fmt.Println("failed to parse the following lines")
		help(false)
		return false
	}

	uuidstr := id
	timestamp := lineQuery[uuidstr].Timestamp
	ip := lineQuery[uuidstr].Stream["_ip_"]
	filename := lineQuery[uuidstr].Stream["filename"]
	taginfo := lineQuery[uuidstr].StreamJsonStr
	threadid := lineQuery[uuidstr].ThreadID
	if traceid != "" || threadid != "" {
		mode = "traceid"
	}
	if timestamp == "" {
		fmt.Println("请先使用more命令定位到关键字当前行的唯一ID: more | grep [keyword] ,再根据唯一ID查询上下文")
		return false
	}
	env := appenvMap[user+":"+app_name]

	filterParam := map[string][]string{}
	filterParam["_app_"] = append(filterParam["_app_"], app_name)
	keyword := "nokey"
	threadglobal = ""
	switch mode {
	case "ip":
		filterParam["_ip_"] = append(filterParam["_ip_"], ip)
	case "filename":
		filterParam["_ip_"] = append(filterParam["_ip_"], ip)
		filterParam["filename"] = append(filterParam["filename"], filename)
		// filelist := logtypeMapGlobal["application"]
		// filetemp := RemoveRepByMap(filelist)
		// filterParam["filename"] = append(filterParam["filename"], filetemp...)
	case "traceid":
		filterParam["_ip_"] = append(filterParam["_ip_"], ip)
		// filterParam["filename"] = append(filterParam["filename"], filename)
		if traceid != "" {
			keyword = traceid
			threadglobal = traceid
		} else if threadid != "" {
			keyword = threadid
			threadglobal = threadid
		}
	}
	// fmt.Println(filterParam)

	// layoutdisplay := "2006-01-02 15:04:05"

	currentTimestampStr := timestamp
	timestampIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
	currentTime := time.Unix(0, timestampIntValue)
	// currentTimeDisplay := currentTime.Format(layoutdisplay)
	// pageTag := currentTimeDisplay + " : " + threadid
	// pageTaghighlight := fmt.Sprintf(blueCode, pageTag)
	// fmt.Println(pageTaghighlight)

	timeunit := ""
	if keyword != "nokey" {
		timeunit = "10m"
	} else {
		timeunit = "1h"
	}
	hsub, _ := time.ParseDuration("-" + timeunit)
	start1time := currentTime.Add(hsub)
	// start1timeDisplay := start1time.Format(layout)
	// fmt.Println(start1timeDisplay)
	start1Timestamp := start1time.UnixNano()
	start1TimestampStr := strconv.FormatInt(start1Timestamp, 10)

	hadd, _ := time.ParseDuration(timeunit)
	end2time := currentTime.Add(hadd)
	// end2timeDisplay := end2time.Format(layout)
	// fmt.Println(end2timeDisplay)
	end2Timestamp := end2time.UnixNano()
	end2TimestampStr := strconv.FormatInt(end2Timestamp, 10)

	contentBrowsePromt := fmt.Sprintf(pageCode, "上下文浏览")
	fmt.Println(contentBrowsePromt)
	switch border {
	case "c":
		backwardres := kernel.QueryRange(filterParam, start1TimestampStr, currentTimestampStr, "5", upnum, keyword, "BACKWARD", "prod", false, user, env)
		doDataProcessing(backwardres, keyword, false, taginfo, "backward", upnumint, true)

		downnumint = downnumint + 1
		downnumstr := strconv.Itoa(downnumint)
		forwardres := kernel.QueryRange(filterParam, currentTimestampStr, end2TimestampStr, "5", downnumstr, keyword, "FORWARD", "prod", true, user, env)
		doDataProcessing(forwardres, keyword, true, taginfo, "forward", 0, true)

		isExistedData := checkIsExistedData(backwardres, forwardres, upnumint, downnumint)
		if isExistedData {
			return true
		}
	case "a":
		downnumint = downnumint + 1
		downnumstr := strconv.Itoa(downnumint)
		response := kernel.QueryRange(filterParam, currentTimestampStr, end2TimestampStr, "5", downnumstr, keyword, "FORWARD", "prod", false, user, env)
		doDataProcessing(response, keyword, true, taginfo, "forward", 0, true)

		isExistedData := checkIsExistedDataSingle(response, downnumint)
		if isExistedData {
			return true
		}
	case "b":
		response := kernel.QueryRange(filterParam, start1TimestampStr, currentTimestampStr, "5", upnum, keyword, "BACKWARD", "prod", false, user, env)
		doDataProcessing(response, keyword, false, taginfo, "backward", upnumint, true)

		isExistedData := checkIsExistedDataSingle(response, upnumint)
		if isExistedData {
			return true
		}
	default:
		break
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()
	promt_msg := "--Context-- Press Space to Paging / Press ESC or CtrlC to quit"
	promt := fmt.Sprintf(pageCode, promt_msg)
	fmt.Println(promt + "\t" + contentBrowsePromt)

	border = "a"
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		if key == keyboard.KeySpace {

			switch border {
			case "c":
				upnumint = upnumint + 20
				upnum = strconv.Itoa(upnumint)
				backwardres := kernel.QueryRange(filterParam, start1TimestampStr, currentTimestampStr, "5", upnum, keyword, "BACKWARD", "prod", false, user, env)
				doDataProcessing(backwardres, keyword, false, taginfo, "backward", upnumint, false)

				downnumint = downnumint + 20
				downnumstr := strconv.Itoa(downnumint)
				forwardres := kernel.QueryRange(filterParam, currentTimestampStr, end2TimestampStr, "5", downnumstr, keyword, "FORWARD", "prod", false, user, env)
				doDataProcessing(forwardres, keyword, true, taginfo, "forward", downnumint-20, false)

				isExistedData := checkIsExistedData(backwardres, forwardres, upnumint, downnumint)
				if isExistedData {
					return true
				}
			case "a":
				if maxLastTimestamp >= end2TimestampStr {
					return true
				}

				downnumint = 21
				downnumstr := strconv.Itoa(downnumint)
				response := kernel.QueryRange(filterParam, maxLastTimestamp, end2TimestampStr, "5", downnumstr, keyword, "FORWARD", "prod", false, user, env)
				doDataProcessing(response, keyword, true, taginfo, "forward", lastrownum, false)

				isExistedData := checkIsExistedDataSingle(response, downnumint)
				if isExistedData {
					return true
				}
			case "b":
				upnumint = upnumint + 20
				upnum = strconv.Itoa(upnumint)
				response := kernel.QueryRange(filterParam, start1TimestampStr, currentTimestampStr, "5", upnum, keyword, "BACKWARD", "prod", false, user, env)
				doDataProcessing(response, keyword, false, taginfo, "backward", upnumint, false)

				isExistedData := checkIsExistedDataSingle(response, upnumint)
				if isExistedData {
					return true
				}
			default:
				break
			}

			promt := fmt.Sprintf(pageCode, promt_msg)
			fmt.Println(promt + "\t" + contentBrowsePromt)

		}
		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			break
		}
	}

	return true
}

/**
 * isGrephighlight:backward无需高亮(false),forward需高亮(true)
 * highlightSpecifiedRow:第一次上下文需要确定高亮行(true),后续翻页无需高亮(false)
 * 小结:只有第一次上下文且forward查询时才高亮关键字所在行
 */
func doDataProcessing(response kernel.QueryRangeGenerate, keyword string, isGrephighlight bool, taginfo string, direction string, limit int, highlightSpecifiedRow bool) {
	keywordHighlight := false
	rownum := true
	events := []map[string]string{}
	line := 0
	switch direction {
	case "forward":
		events = doLogProcessingContext(response, keyword, keywordHighlight, rownum, limit)
	case "backward":
		events = doLogProcessingContextBackward(response, keyword, keywordHighlight, rownum, limit)
		totalrows := countTotalRows(response)
		// fmt.Println(totalrows)
		line = ^totalrows + 1
	}

	// 对 events 按照timestamp进行排序
	events = mapSort(events, "timestampNano", "forward")

	if highlightSpecifiedRow {
		displayContextData(events, isGrephighlight, taginfo, direction, line)
	} else {
		displayContextDataPageing(events, isGrephighlight, taginfo)
	}

	// showFilterCondition()
}

func displayContextData(events []map[string]string, isGrephighlight bool, taginfo string, direction string, line int) {
	highlightSpecifiedRow := true

	for _, row := range events {

		// tagDisplay := fmt.Sprintf(blueCode, row["tag"])
		// timeDisplay := fmt.Sprintf(infoCode, row["timestamp"])

		if row["tag"] == taginfo && isGrephighlight && highlightSpecifiedRow {
			row["content"] = fmt.Sprintf(warningBoldCode, row["content"])
			highlightSpecifiedRow = false
		} else {
			if threadglobal != "" {
				hightThreadid := fmt.Sprintf(warningCode, threadglobal)
				row["content"] = strings.Replace(row["content"], threadglobal, hightThreadid, 1)
			}
		}

		labelDisplay := addRownumAndTimeToContextTag(row["timestamp"], direction, line, row["tag"])
		line++
		contentDisplay := addMarkToContextContent(row["content"])

		if row["content"] != "" {
			fmt.Println(labelDisplay)
			fmt.Println(contentDisplay)
			fmt.Println("-----------------------------------------------------------------------------------------------------------------------------------------------")
		}

		// laststarttime = row["timestamp"]
	}
}

func displayContextDataPageing(events []map[string]string, isGrephighlight bool, taginfo string) {
	isPreviousPageLastLineShow := true
	line := lastrownum

	for _, row := range events {

		// tagDisplay := fmt.Sprintf(whiteCode, row["tag"])
		// timeDisplay := fmt.Sprintf(infoCode, row["timestamp"])

		if row["tag"] == lasttaginfo && isGrephighlight && isPreviousPageLastLineShow {
			isPreviousPageLastLineShow = false
			line++
			continue
		}
		if threadglobal != "" {
			hightThreadid := fmt.Sprintf(warningCode, threadglobal)
			row["content"] = strings.Replace(row["content"], threadglobal, hightThreadid, 1)
		}

		labelDisplay := addRownumAndTimeToContextTag(row["timestamp"], "forward", line, row["tag"])
		line++
		contentDisplay := addMarkToContextContent(row["content"])

		if row["content"] != "" {
			fmt.Println(labelDisplay)
			fmt.Println(contentDisplay)
			fmt.Println("-----------------------------------------------------------------------------------------------------------------------------------------------")
		}
	}
}

func addRownumAndTimeToContextTag(timestamp string, direction string, line int, tagSource string) string {
	if direction == "forward" {
		lastrownum = line
		lasttaginfo = tagSource
	}

	linenum := strconv.Itoa(line)
	hllinenum := fmt.Sprintf(redCode, linenum)
	if line == 0 {
		hllinenum = fmt.Sprintf(successCode, linenum)
	}
	if line > 0 {
		hllinenum = fmt.Sprintf(successCode, "+"+linenum)
	}

	hltime := fmt.Sprintf(dateCode, "["+timestamp+"]")
	tagDisplay := fmt.Sprintf(blueCode, tagSource)
	labelDisplay := hllinenum + " " + hltime + tagDisplay

	return labelDisplay
}

func addMarkToContextContent(content string) string {
	// contentMark := fmt.Sprintf(contentMarkCode, "content")
	contentMark := "\t" + "content:"
	content = contentMark + "\n\t" + content
	return content
}

func checkIsExistedData(backwardres kernel.QueryRangeGenerate, forwardres kernel.QueryRangeGenerate, backwardlimit int, forwardlimit int) bool {
	backwardrownum := searchLastTimeAndMaxRowNum(backwardres)
	forwardrownum := searchLastTimeAndMaxRowNum(forwardres)
	// fmt.Println(backwardrownum)
	// fmt.Println(forwardrownum)
	maxLastTimestamp = forwardrownum["maxLastTimestamp"].(string)

	if backwardrownum["aliyuntotalrownum"].(int) < backwardlimit &&
		backwardrownum["upyuntotalrownum"].(int) < backwardlimit &&
		backwardrownum["totalnum"].(int) < backwardlimit &&
		forwardrownum["aliyuntotalrownum"].(int) < forwardlimit &&
		forwardrownum["upyuntotalrownum"].(int) < forwardlimit &&
		forwardrownum["totalnum"].(int) < forwardlimit {
		fmt.Println("该时间段内已没有上下文数据")
		return true
	}
	return false
}

func checkIsExistedDataSingle(response kernel.QueryRangeGenerate, limit int) bool {
	responserownum := searchLastTimeAndMaxRowNum(response)
	// fmt.Println(responserownum)
	maxLastTimestamp = responserownum["maxLastTimestamp"].(string)

	if responserownum["aliyuntotalrownum"].(int) < limit &&
		responserownum["upyuntotalrownum"].(int) < limit && responserownum["totalnum"].(int) < limit {
		fmt.Println("该时间段内已没有上下文数据")
		return true
	}
	return false
}

func displayFinalDataPageingMultiLineJoin(events []map[string]string, isGrephighlight bool, taginfo string) {
	highlightSpecifiedRow := true
	isPreviousPageLastLineShow := false
	line := lastrownum

	for _, row := range events {

		tagDisplay := fmt.Sprintf(blueCode, row["tag"])
		// timeDisplay := fmt.Sprintf(infoCode, row["timestamp"])
		// timestampNano := fmt.Sprintf(infoCode, row["timestampNano"])

		if row["tag"] == lasttaginfo && isGrephighlight && highlightSpecifiedRow {
			timeContent := strings.Split(row["content"], "\n")
			for index, content := range timeContent {
				//省略显示翻页后紧接着上一页最后一行的内容
				if index == 0 {
					continue
					// backgroundLine := fmt.Sprintf(whiteCode, content)
					// contentline := strings.Replace(content, content, backgroundLine, 1)
					// row["content"] = contentline
				} else {
					hightThreadid := fmt.Sprintf(warningCode, threadglobal)
					contentline := strings.Replace(content, threadglobal, hightThreadid, 1)
					row["content"] = row["content"] + "\n" + contentline

					isPreviousPageLastLineShow = true
				}
			}

			highlightSpecifiedRow = false
		} else {
			hightThreadid := fmt.Sprintf(warningCode, threadglobal)
			row["content"] = strings.Replace(row["content"], threadglobal, hightThreadid, -1)

			isPreviousPageLastLineShow = true
		}

		if isGrephighlight {
			lastrownum = line
			lasttaginfo = row["tag"]

			linenum := strconv.Itoa(line)
			hightlightlinenum := fmt.Sprintf(purpleCode, linenum)
			row["content"] = hightlightlinenum + "\t" + row["content"]
			line++
		}

		if !highlightSpecifiedRow && !isPreviousPageLastLineShow {
			continue
		}

		if row["content"] != "" {
			fmt.Println(tagDisplay)
			// fmt.Println(timeDisplay + " " + row["content"])
			fmt.Println(row["content"])
		}

		// laststarttime = row["timestamp"]
	}
}

func showFilterCondition() {
	fmt.Println("-----Filter--------")
	fmt.Println("StartTime:" + starttime)
	fmt.Println("EndTime:" + endtime)
	// fmt.Println("StarTimeUnix:" + strconv.FormatInt(startTime, 10))
	// fmt.Println("EndTimeUnix:" + strconv.FormatInt(endTime, 10))
	fmt.Println("Filter:" + util.MapToJson(filterMap))
	// fmt.Println("Step:" + step)
	// fmt.Println("Limit:" + limit)
	// fmt.Println("Direction:" + direction)
	// fmt.Println("Grep:" + keyword)
	fmt.Println("-----------------------------")
}
