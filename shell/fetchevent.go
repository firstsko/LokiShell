package shell

import (
	"encoding/json"
	"fmt"
	"lokishell/kernel"
	"lokishell/util"
	"sort"
	"strconv"
	"strings"
	"time"
)

var starttime = ""
var endtime = ""
var lastmoreline = 1
var lasttimestamp = ""

// var eventsBackup = []map[string]string{}

var eventsBackup = []string{}
var lineQuery = map[string]QueryCondition{}

var tagsAssociatedMap = map[string]int{
	"app":               1,
	"_app_":             1,
	"_appid_":           1,
	"_service_name_":    1,
	"_node_name_":       1,
	"_pod_uid_":         1,
	"_container_name_":  1,
	"stream":            1,
	"_namespace_":       1,
	"_k8s_app_":         1,
	"_pod_name_":        1,
	"branch":            1,
	"deploymentID":      1,
	"group":             1,
	"groupid":           1,
	"mesh_mode":         1,
	"mesh_version":      1,
	"pod_template_hash": 1,
	"region":            1,
}

type QueryCondition struct {
	Stream        map[string]string `json:"stream"`
	Timestamp     string            `json:"timestamp"`
	StreamJsonStr string            `json:"streamJsonStr"`
	ThreadID      string            `json:"threadID"`
}

type MapsSort struct {
	MapList []map[string]string
	by      func(p, q map[string]string) bool
}
type SortBy func(p, q map[string]string) bool

func (m *MapsSort) Len() int {
	return len(m.MapList)
}

func (m *MapsSort) Less(i, j int) bool {
	return m.by(m.MapList[i], m.MapList[j])
}

func (m *MapsSort) Swap(i, j int) {
	m.MapList[i], m.MapList[j] = m.MapList[j], m.MapList[i]
}

func mapSort(ms []map[string]string, key string, direction string) []map[string]string {
	mapsSort := MapsSort{}
	mapsSort.by = func(p, q map[string]string) bool {
		value := p[key] < q[key]
		if direction == "BACKWARD" {
			value = p[key] > q[key]
		}
		return value
	}
	mapsSort.MapList = ms
	sort.Sort(&mapsSort)
	return mapsSort.MapList
}

func fetchEvent(filterMap map[string][]string, starttime string, endtime string, step string, limit string, keyword string, direction string, displayFilter bool, isMorePaging bool) map[string]interface{} {

	// startTimeUnix := strconv.Itoa(int(startTime))
	// endTimeUnix := strconv.Itoa(int(endTime))
	env := appenvMap[user+":"+app_name]
	// fmt.Println("env: " + env)

	var response kernel.QueryRangeGenerate
	if !isMorePaging {
		response = kernel.QueryRange(filterMap, starttime, endtime, step, limit, keyword, direction, "prod", false, user, env)
	} else {
		limitint, _ := strconv.Atoi(limit)
		limitint++
		limitstr := strconv.Itoa(limitint)
		if direction == "BACKWARD" {
			limitstr = limit
		}
		response = kernel.QueryRange(filterMap, starttime, endtime, step, limitstr, keyword, direction, "prod", false, user, env)
	}

	forwardrownum := searchLastTimeAndMaxRowNum(response)

	events := doLogProcessingOriginal(response, keyword, true, false)
	// events = doLogProcessingV2(response, keyword, true, false, 0)

	events = mapSort(events, "timestampNano", direction)
	// eventsBackup = events

	displayMoreQueryResult(events, isMorePaging)

	displayFilter = false
	if displayFilter {
		showMoreQueryFilterInfo(starttime, endtime, limit, keyword)
	}

	return forwardrownum
}

func displayMoreQueryResult(events []map[string]string, isMorePaging bool) {
	moreline := 0
	if isMorePaging {
		moreline = lastmoreline
	} else {
		eventsBackup = []string{}
	}
	for index, row := range events {
		if isMorePaging && index == 0 {
			timestamp := row["timestampNano"]
			if timestamp == lasttimestamp {
				continue
			}
		}
		moreline++
		doAddInfoToResultAndDisplay(row, moreline)
	}
}

func doAddInfoToResultAndDisplay(row map[string]string, linenum int) {
	// tagDisplay := fmt.Sprintf(blueCode, row["tag"])
	// timeDisplay := fmt.Sprintf(infoCode, row["timestamp"])

	labelDisplay := addRownumAndDateToMoreTag(row["tag"], row["timestamp"], linenum, row["timestampNano"])
	contentDisplay := addMarkToMoreContent(row["content"])

	if row["content"] != "" {
		showQueryReslult(labelDisplay, contentDisplay)
	}
}

func showQueryReslult(labelDisplay string, contentDisplay string) {
	fmt.Println(labelDisplay)
	// fmt.Println(timeDisplay + " " + row["content"])
	fmt.Println(contentDisplay)
	dividingLine := "-----------------------------------------------------------------------------------------------------------------------------------------------"
	fmt.Println(dividingLine)
	eventsBackup = append(eventsBackup, labelDisplay)
	eventsBackup = append(eventsBackup, contentDisplay)
	eventsBackup = append(eventsBackup, dividingLine)
}

func displayRecentOnceQueryResult() {
	for _, value := range eventsBackup {
		fmt.Println(value)
	}
}

func addRownumAndDateToMoreTag(tag string, timestamp string, line int, timestampnano string) string {
	linestr := strconv.Itoa(line)
	hlline := fmt.Sprintf(blueCode, linestr)
	hldate := fmt.Sprintf(dateCode, timestamp)
	tagDisplay := fmt.Sprintf(blueCode, tag)
	labelDisplay := hlline + hldate + tagDisplay

	lastmoreline = line
	lasttimestamp = timestampnano
	return labelDisplay
}

func addMarkToMoreContent(content string) string {
	contentMark := fmt.Sprintf(contentMarkCode, "content")
	// content = "\t\t    " + contentMark + ":" + content
	content = "\t   " + contentMark + ":\n\t   " + content
	return content
}

func showMoreQueryFilterInfo(starttime string, endtime string, limit string, keyword string) {
	starttimeFormat := util.TimestampNanoStrToFormattime(starttime)
	endtimeFormat := util.TimestampNanoStrToFormattime(endtime)

	filterInfo := map[string]string{
		"StartTime":    starttimeFormat,
		"EndTime":      endtimeFormat,
		"StarTimeUnix": starttime,
		"EndTimeUnix":  endtime,
		"Filter":       util.MapToJson(filterMap),
		"Limit":        limit,
		"Grep":         keyword,
		// "Step":         step,
		// "Direction":    direction,
	}

	// fmt.Println("--------------------------FilterInfo--------------------------")
	fmt.Println()
	fmt.Print("FilterInfo: ")
	fmt.Println(filterInfo)
	fmt.Println()
}

func doLogProcessingOriginal(response kernel.QueryRangeGenerate, keyword string, keywordHighlight bool, rownum bool) []map[string]string {
	eventlist := []map[string]string{}

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {
		line := 0
		for _, stream := range response.Data.Data.Result {

			tags := excludeAssociatedTags(stream.Stream)

			for _, value := range stream.Values {
				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("01-02 15:04:05")

				tagsInfo := tags
				if isHasKeywordGrep {
					tagsInfo = addUUIDAndAnalysisContent(value[1].(string), stream.Stream, tags, timestamp, "tag")
				}

				// 高亮关键字：区别关键字之间且与或关系
				highlightcontent := value[1].(string)
				if keyword != "" && keyword != "nokey" && keywordHighlight {
					highlightcontent = highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)
				}

				if rownum {
					linenum := strconv.Itoa(line)
					highlightcontent = linenum + "\t" + highlightcontent
				}

				data := map[string]string{
					"tag":           tagsInfo,
					"timestamp":     timeDisplay,
					"content":       highlightcontent,
					"timestampNano": timestamp,
				}

				eventlist = append(eventlist, data)

			}
		}

	} else {

		// errCode := strconv.Itoa(response.Data.RetCode)
		// errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)
		// fmt.Println(errMsg)

		promtMsg := fmt.Sprintf(warningCode, "当前查询条件无相关结果")
		fmt.Println(promtMsg)
	}

	return eventlist
}

func doLogProcessingContext(response kernel.QueryRangeGenerate, keyword string, keywordHighlight bool, rownum bool, limit int) []map[string]string {
	eventlist := []map[string]string{}

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {

		for _, stream := range response.Data.Data.Result {

			tags := excludeAssociatedTags(stream.Stream)

			for _, value := range stream.Values {
				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("2006-01-02 15:04:05")

				// 高亮关键字：区别关键字之间且与或关系
				content := value[1].(string)
				if keyword != "" && keyword != "nokey" && keywordHighlight {
					content = highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)
				}

				data := map[string]string{
					"tag":           tags,
					"timestamp":     timeDisplay,
					"content":       content,
					"timestampNano": timestamp,
				}

				eventlist = append(eventlist, data)
			}
		}
	} else {

		// errCode := strconv.Itoa(response.Data.Code)
		// errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)
		// fmt.Println(errMsg)

		promtMsg := fmt.Sprintf(warningCode, "当前查询无下文结果")
		fmt.Println(promtMsg)
	}

	return eventlist
}

func doLogProcessingV2(response kernel.QueryRangeGenerate, keyword string, keywordHighlight bool, rownum bool, limit int) []map[string]string {
	eventlist := []map[string]string{}
	contents := strings.Builder{}
	tagsInit := ""
	timeDisplayInit := ""

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {
		line := limit
		for i, stream := range response.Data.Data.Result {

			tags := excludeAssociatedTags(stream.Stream)
			if i == 0 {
				tagsInit = tags
			}

			for index, value := range stream.Values {
				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("2006-01-02 15:04:05")
				if i == 0 && index == 0 {
					timeDisplayInit = timeDisplay
				}

				if tagsInit != tags || timeDisplayInit != timeDisplay {
					data := map[string]string{
						"tag":       tagsInit,
						"timestamp": timeDisplayInit,
						"content":   contents.String(),
					}

					eventlist = append(eventlist, data)

					if tagsInit != tags {
						tagsInit = tags
					}
					if timeDisplayInit != timeDisplay {
						timeDisplayInit = timeDisplay
					}

					contents.Reset()
				}

				highlightcontent := value[1].(string)
				if isHasKeywordGrep {
					highlightcontent = addUUIDAndAnalysisContent(highlightcontent, stream.Stream, tags, timestamp, "content")
				}
				// 高亮关键字：区别关键字之间且与或关系
				if keyword != "" && keyword != "nokey" && keywordHighlight {
					highlightcontent = highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)
				}
				// 加上行号
				if rownum {
					lastrownum = line

					linenum := strconv.Itoa(line)
					hightlightlinenum := fmt.Sprintf(purpleCode, linenum)
					contents.WriteString(hightlightlinenum + "\t")
					line++
				}

				contents.WriteString(highlightcontent)
				contents.WriteString("\n")
			}
		}

		// 添加最后一条记录
		if timeDisplayInit != "" {
			data := map[string]string{
				"tag":       tagsInit,
				"timestamp": timeDisplayInit,
				"content":   contents.String(),
			}

			eventlist = append(eventlist, data)
		}

	} else {

		errCode := strconv.Itoa(response.Data.RetCode)
		errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)

		fmt.Println(errMsg)
	}

	return eventlist
}

func excludeAssociatedTags(tagMap map[string]string) string {
	for key := range tagMap {
		if value, ok := tagsAssociatedMap[key]; ok {
			if value == 1 {
				delete(tagMap, key)
			}
		}
	}
	bytes, _ := json.Marshal(tagMap)
	tags := string(bytes)

	return tags
}

func addUUIDAndAnalysisContent(content string, stream map[string]string, tags string, timestamp string, idplace string) string {
	// uid, err := uuid.NewV4()
	// if err != nil {
	// 	fmt.Printf("Something went wrong with uuid")
	// 	return ""
	// }
	// uidstr := uid.String()
	threaduid := util.CharacterRandomNumber(12)
	hlthreaduid := fmt.Sprintf(redCode, threaduid)

	fileuid := util.CharacterRandomNumber(12)
	hlfileuid := fmt.Sprintf(redCode, fileuid)

	newInfo := ""
	switch idplace {
	case "tag":
		FilehlPrompt := fmt.Sprintf(blueCode, "日志上下文")
		ThreadhlPrompt := fmt.Sprintf(blueCode, "线程上下文")
		newInfo = tags + "\t" + FilehlPrompt + ":" + hlfileuid + " " + ThreadhlPrompt + ": " + hlthreaduid
	case "content":
		newInfo = content + "\t" + hlthreaduid
	}

	encodingMappingContent(content, stream, tags, timestamp, threaduid, true)

	encodingMappingContent("", stream, tags, timestamp, fileuid, false)

	return newInfo
}

func encodingMappingContent(content string, stream map[string]string, tags string, timestamp string, uid string, isThreadContext bool) {
	threadID := ""
	if isThreadContext {
		threadID = analysisThreadId(content)
	}

	mappingct := new(QueryCondition)
	mappingct.Stream = stream
	mappingct.Timestamp = timestamp
	mappingct.StreamJsonStr = tags
	mappingct.ThreadID = threadID
	lineQuery[uid] = *mappingct
}

func analysisThreadId(content string) string {
	threadid := ""
	contentarr := strings.Fields(content)
	// FATAL/ERROR/WARN/INFO/DEBUG/TRACE
	isExistedLogLevel := false
	for index, value := range contentarr {
		switch value {
		case "FATAL":
			fallthrough
		case "ERROR":
			fallthrough
		case "WARN":
			fallthrough
		case "INFO":
			fallthrough
		case "DEBUG":
			fallthrough
		case "TRACE":
			threadid = parseThread(index, contentarr)
			isExistedLogLevel = true
		}
		if isExistedLogLevel {
			break
		}
	}

	return threadid
}

func parseThread(loglevelIndex int, contentarr []string) string {
	threadid := ""
	threadidindex := loglevelIndex + 1
	for index, value := range contentarr {
		if index > threadidindex {
			threadvalue := value
			endFlag := false
			if strings.HasSuffix(value, "]") {
				endIndex := strings.Index(value, "]")
				threadvalue = value[:endIndex+1]
				endFlag = true
			}
			threadid += " " + threadvalue
			if endFlag {
				break
			}
		}
		if index == threadidindex {
			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				endIndex := strings.Index(value, "]")
				threadid = value[:endIndex+1]
				break
			} else if strings.HasPrefix(value, "[") {
				threadid = value
				continue
			}
		}
	}
	// fmt.Println(threadid)
	return threadid
}

func doLogProcessingOriginalBackward(response kernel.QueryRangeGenerate, keyword string, keywordHighlight bool, rownum bool) []map[string]string {
	eventlist := []map[string]string{}

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {
		line := 0

		result := response.Data.Data.Result
		resultlen := len(result)
		for i := resultlen - 1; i >= 0; i-- {
			tags := excludeAssociatedTags(result[i].Stream)

			values := result[i].Values
			valueslen := len(values)
			for index := valueslen - 1; index >= 0; index-- {
				value := values[index]
				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("2006-01-02 15:04:05")

				tagsInfo := tags
				if isHasKeywordGrep {
					tagsInfo = addUUIDAndAnalysisContent(value[1].(string), result[i].Stream, tags, timestamp, "tag")
				}

				// 高亮关键字：区别关键字之间且与或关系
				highlightcontent := value[1].(string)
				if keyword != "" && keyword != "nokey" && keywordHighlight {
					highlightcontent = highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)
				}

				if rownum {
					linenum := strconv.Itoa(line)
					highlightcontent = linenum + "\t" + highlightcontent
				}

				data := map[string]string{
					"tag":           tagsInfo,
					"timestamp":     timeDisplay,
					"content":       highlightcontent,
					"timestampNano": timestamp,
				}

				eventlist = append(eventlist, data)
			}
		}
	} else {

		// errCode := strconv.Itoa(response.Data.RetCode)
		// errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)
		// fmt.Println(errMsg)

		promtMsg := fmt.Sprintf(warningCode, "当前查询条件无相关结果")
		fmt.Println(promtMsg)
	}

	return eventlist
}

func doLogProcessingContextBackward(response kernel.QueryRangeGenerate, keyword string, keywordHighlight bool, rownum bool, limit int) []map[string]string {
	eventlist := []map[string]string{}

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {
		// totalrows := countTotalRows(response)
		// line := ^totalrows + 1

		result := response.Data.Data.Result
		resultlen := len(result)
		for i := resultlen - 1; i >= 0; i-- {
			tags := excludeAssociatedTags(result[i].Stream)

			values := result[i].Values
			valueslen := len(values)
			for index := valueslen - 1; index >= 0; index-- {
				value := values[index]
				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("2006-01-02 15:04:05")

				// 高亮关键字：区别关键字之间且与或关系
				content := value[1].(string)
				if keyword != "" && keyword != "nokey" && keywordHighlight {
					content = highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)
				}

				// if rownum {
				// 	linenum := strconv.Itoa(line)
				// 	hightlightlinenum := fmt.Sprintf(purpleCode, linenum)
				// 	content = hightlightlinenum + "\t" + content
				// 	line++
				// }

				data := map[string]string{
					"tag":           tags,
					"timestamp":     timeDisplay,
					"content":       content,
					"timestampNano": timestamp,
				}

				eventlist = append(eventlist, data)
			}
		}
	} else {

		// errCode := strconv.Itoa(response.Data.Code)
		// errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)

		promtMsg := fmt.Sprintf(warningCode, "当前查询无上文结果")
		fmt.Println(promtMsg)
	}

	return eventlist
}

func doLogProcessingV2Backward(response kernel.QueryRangeGenerate, keyword string, keywordHighlight bool, rownum bool, limit int) []map[string]string {
	eventlist := []map[string]string{}
	contents := strings.Builder{}
	tagsInit := ""
	timeDisplayInit := ""

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {
		totalrows := countTotalRows(response)
		line := ^totalrows + 1

		result := response.Data.Data.Result
		resultlen := len(result)
		for i := resultlen - 1; i >= 0; i-- {
			tags := excludeAssociatedTags(result[i].Stream)
			if i == 0 {
				tagsInit = tags
			}

			values := result[i].Values
			valueslen := len(values)
			for index := valueslen - 1; index >= 0; index-- {
				value := values[index]
				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("2006-01-02 15:04:05")
				if i == 0 && index == 0 {
					timeDisplayInit = timeDisplay
				}

				if tagsInit != tags || timeDisplayInit != timeDisplay {
					data := map[string]string{
						"tag":       tagsInit,
						"timestamp": timeDisplayInit,
						"content":   contents.String(),
					}

					eventlist = append(eventlist, data)

					if tagsInit != tags {
						tagsInit = tags
					}
					if timeDisplayInit != timeDisplay {
						timeDisplayInit = timeDisplay
					}

					contents.Reset()
				}

				highlightcontent := value[1].(string)
				if isHasKeywordGrep {
					highlightcontent = addUUIDAndAnalysisContent(highlightcontent, result[i].Stream, tags, timestamp, "content")
				}
				// 高亮关键字：区别关键字之间且与或关系
				if keyword != "" && keyword != "nokey" && keywordHighlight {
					highlightcontent = highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)
				}

				if rownum {
					linenum := strconv.Itoa(line)
					hightlightlinenum := fmt.Sprintf(purpleCode, linenum)
					contents.WriteString(hightlightlinenum + "\t")
					line++
				}

				contents.WriteString(highlightcontent)
				contents.WriteString("\n")

			}

		}

		// 添加最后一条记录
		if timeDisplayInit != "" {
			data := map[string]string{
				"tag":       tagsInit,
				"timestamp": timeDisplayInit,
				"content":   contents.String(),
			}

			eventlist = append(eventlist, data)
		}

	} else {

		errCode := strconv.Itoa(response.Data.RetCode)
		errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)

		fmt.Println(errMsg)
	}

	return eventlist
}

func countTotalRows(response kernel.QueryRangeGenerate) int {
	aliyuntotalrownum := 0
	upyuntotalrownum := 0
	totalnum := 0
	if len(response.Data.Data.Result) > 0 {
		for _, stream := range response.Data.Data.Result {
			rowcount := len(stream.Values)
			totalnum += rowcount

			ip := stream.Stream["_ip_"]
			switch ipidcMap[ip] {
			case "aliyun":
				aliyuntotalrownum += rowcount
			case "upyun":
				upyuntotalrownum += rowcount
			}
		}
	}
	total := aliyuntotalrownum + upyuntotalrownum
	if total == 0 {
		total = totalnum
	}
	return total
}

func searchLastTimeAndMaxRowNum(response kernel.QueryRangeGenerate) map[string]interface{} {
	rowlist := []int{}
	lasttimestamplist := []string{}
	firsttimestamplist := []string{}
	// lasttimelist := []string{}
	// firsttimelist := []string{}
	aliyuntotalrownum := 0
	upyuntotalrownum := 0
	totalnum := 0
	if len(response.Data.Data.Result) > 0 {
		for _, stream := range response.Data.Data.Result {
			rowcount := len(stream.Values)
			ip := stream.Stream["_ip_"]
			switch ipidcMap[ip] {
			case "aliyun":
				aliyuntotalrownum += rowcount
			case "upyun":
				upyuntotalrownum += rowcount
			}
			totalnum += rowcount
			rowlist = append(rowlist, rowcount)
			if rowcount > 0 {
				firsttimestamp := stream.Values[0][0].(string)
				lasttimestamp := stream.Values[rowcount-1][0].(string)
				// firsttime := timestampFormat(firsttimestamp)
				// lasttime := timestampFormat(lasttimestamp)
				if firsttimestamp < lasttimestamp {
					firsttimestamplist = append(firsttimestamplist, firsttimestamp)
					lasttimestamplist = append(lasttimestamplist, lasttimestamp)
				} else {
					firsttimestamplist = append(firsttimestamplist, lasttimestamp)
					lasttimestamplist = append(lasttimestamplist, firsttimestamp)
				}
				// if firsttime < lasttime {
				// 	lasttimelist = append(lasttimelist, lasttime)
				// 	firsttimelist = append(firsttimelist, firsttime)
				// } else {
				// 	lasttimelist = append(lasttimelist, firsttime)
				// 	firsttimelist = append(firsttimelist, lasttime)
				// }
			}
		}
	}

	maxRowNum := getMinOrMaxOfIntList(rowlist, "max")
	// minLastTime := getMinOrMaxOfStringList(lasttimelist, "min")
	// minFirstTime := getMinOrMaxOfStringList(firsttimelist, "min")
	maxLastTimestamp := getMinOrMaxOfStringList(lasttimestamplist, "max")
	minFirstTimestamp := getMinOrMaxOfStringList(firsttimestamplist, "min")

	controlParams = map[string]interface{}{
		// "lastrownum":        lastrownum,
		"minFirstTimestamp": minFirstTimestamp,
		"maxLastTimestamp":  maxLastTimestamp,
		// "firststarttime":    minFirstTime,
		// "laststarttime":     minLastTime,
		"maxRowNum":         maxRowNum,
		"aliyuntotalrownum": aliyuntotalrownum,
		"upyuntotalrownum":  upyuntotalrownum,
		"totalnum":          totalnum,
	}

	return controlParams
}

func getMinOrMaxOfIntList(list []int, valuetype string) int {
	value := 0
	if len(list) > 0 {
		if !sort.IntsAreSorted(list) {
			sort.Ints(list)
		}
		if valuetype == "max" {
			value = list[len(list)-1]
		} else {
			value = list[0]
		}
	}
	return value
}

func getMinOrMaxOfStringList(list []string, valuetype string) string {
	value := ""
	if len(list) > 0 {
		if !sort.StringsAreSorted(list) {
			sort.Strings(list)
		}
		if valuetype == "max" {
			value = list[len(list)-1]
		} else {
			value = list[0]
		}
	}
	return value
}

func timestampFormat(timestamp string) string {
	unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
	t := time.Unix(0, unixIntValue)
	timeDisplay := t.Format("2006-01-02 15:04:05")
	return timeDisplay
}

func highlightKeyword(content string, keywordlist []string, grepAndOrRel string) string {
	replace := map[string]string{}

	for _, value := range keywordlist {
		hightlightKey := fmt.Sprintf(warningCode, value)
		replace[value] = hightlightKey
	}

	isAndRel := true
	if grepAndOrRel == "and" {
		isAndRel = checkAndRelIsMeeted(content, keywordlist)
	}
	if isAndRel {
		for s, r := range replace {
			content = strings.Replace(content, s, r, -1)
		}
	}

	return content
}

func checkAndRelIsMeeted(contents string, keywordlist []string) bool {
	isAndRel := true
	for _, value := range keywordlist {
		if !strings.Contains(contents, value) {
			isAndRel = false
			break
		}
	}
	return isAndRel
}

func doLogProcessing(response kernel.QueryRangeGenerate, keyword string, isGrepContext bool, grepContextParam map[string]map[string]int) []map[string]string {
	eventlist := []map[string]string{}
	contents := strings.Builder{}
	tagsInit := ""
	timeDisplayInit := ""

	grepAndOrRel := "and"
	keywordlist := strings.Split(keyword, ",")
	if len(response.Data.Data.Result) > 0 {

		isShowBitMap := map[int][]bool{}
		if isGrepContext {
			isShowBitMap = showAndHideBitList(response, grepContextParam, keywordlist, grepAndOrRel)
		}

		for i, stream := range response.Data.Data.Result {

			bytes, _ := json.Marshal(stream.Stream)
			tags := string(bytes)
			if i == 0 {
				tagsInit = tags
			}

			isShowBit, ok := isShowBitMap[i]
			showSign := 0
			for index, value := range stream.Values {
				if ok && !isShowBit[index] {
					if showSign == 1 {
						contents.WriteString("---")
						contents.WriteString("\n")
						showSign = 0
					}
					continue
				}
				showSign = 1

				timestamp := value[0].(string)
				unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
				t := time.Unix(0, unixIntValue)
				timeDisplay := t.Format("2006-01-02 15:04:05")
				if index == 0 {
					timeDisplayInit = timeDisplay
				}

				// 高亮关键字：区别关键字之间且与或关系
				highlightcontent := highlightKeyword(value[1].(string), keywordlist, grepAndOrRel)

				if tagsInit != tags || timeDisplayInit != timeDisplay {
					data := map[string]string{
						"tag":       tagsInit,
						"timestamp": timeDisplayInit,
						"content":   contents.String(),
					}

					eventlist = append(eventlist, data)

					if tagsInit != tags {
						tagsInit = tags
					}
					if timeDisplayInit != timeDisplay {
						timeDisplayInit = timeDisplay
					}

					contents.Reset()
				}

				contents.WriteString(highlightcontent)
				contents.WriteString("\n")
			}
		}

		// 添加最后一条记录
		if timeDisplayInit != "" {
			data := map[string]string{
				"tag":       tagsInit,
				"timestamp": timeDisplayInit,
				"content":   contents.String(),
			}

			eventlist = append(eventlist, data)
		}

	} else {

		errCode := strconv.Itoa(response.Data.RetCode)
		errMsg := fmt.Sprintf(warningCode, errCode+" "+response.Data.Message)

		fmt.Println(errMsg)
	}

	return eventlist
}

func showAndHideBitList(response kernel.QueryRangeGenerate, grepContextParam map[string]map[string]int, keywordlist []string, grepAndOrRel string) map[int][]bool {
	isShowBitMap := map[int][]bool{}

	for key, stream := range response.Data.Data.Result {

		rangelist := [][]int{}
		for index, value := range stream.Values {
			contents := value[1].(string)

			for param, mapList := range grepContextParam {
				switch param {
				case "C":
					for keyword, num := range mapList {
						if strings.Contains(contents, keyword) {
							isAndRel := true
							if grepAndOrRel == "and" {
								isAndRel = checkAndRelIsMeeted(contents, keywordlist)
							}
							if isAndRel {
								start := index - num
								end := index + num
								scope := []int{start, end}
								rangelist = append(rangelist, scope)
							}
						}
					}
				case "B":
					for keyword, num := range mapList {
						if strings.Contains(contents, keyword) {
							isAndRel := true
							if grepAndOrRel == "and" {
								isAndRel = checkAndRelIsMeeted(contents, keywordlist)
							}
							if isAndRel {
								start := index - num
								end := index
								scope := []int{start, end}
								rangelist = append(rangelist, scope)
							}
						}
					}
				case "A":
					for keyword, num := range mapList {
						if strings.Contains(contents, keyword) {
							isAndRel := true
							if grepAndOrRel == "and" {
								isAndRel = checkAndRelIsMeeted(contents, keywordlist)
							}
							if isAndRel {
								start := index
								end := index + num
								scope := []int{start, end}
								rangelist = append(rangelist, scope)
							}
						}
					}
				case "D":
					for keyword, _ := range mapList {
						if strings.Contains(contents, keyword) {
							isAndRel := true
							if grepAndOrRel == "and" {
								isAndRel = checkAndRelIsMeeted(contents, keywordlist)
							}
							if isAndRel {
								scope := []int{index, index}
								rangelist = append(rangelist, scope)
							}
						}
					}
				}
			}
		}

		valuesLength := len(stream.Values)
		isShowBit := make([]bool, valuesLength)

		for _, value := range rangelist {
			start := value[0]
			end := value[1]
			for index := start; index <= end; index++ {
				if index < 0 {
					index = 0
				}
				if index > valuesLength-1 {
					index = end
				}
				if index < valuesLength {
					isShowBit[index] = true
				}
			}
		}

		isShowBitMap[key] = isShowBit
	}

	return isShowBitMap
}
