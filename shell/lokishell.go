package shell

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"lokishell/kernel"
	"lokishell/util"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

const redCode = "\033[31m %s \033[0m"
const warningCode = "\033[33m %s \033[0m"
const successCode = "\033[32m %s \033[0m"
const purpleCode = "\033[35m %s \033[0m"
const infoCode = "\033[36m %s \033[0m"
const blueCode = "\033[34m %s \033[0m"
const whiteCode = "\033[37m %s \033[0m"
const warningBoldCode = "\033[1;93m %s \033[0m"
const dateCode = "\033[1;94m %s \033[0m"
const timeCode = "\033[1;97m %s \033[0m"
const pageCode = "\033[1;30;107m %s \033[0m"
const contentMarkCode = "\033[1;30;107m %s \033[0m"
const divideLineCode = "\033[1;30;107m %s \033[0m"
const logcli = "./logcli"

var s_path = "/"
var t_path = "/"
var user = "anonymous"
var token = ""
var prefix = "[" + user + "@lokishell " + s_path + "]$"
var app_name = "hoc-openrestry-ser"
var nodeType = "root"
var layout = "2006-01-0215:04:05"
var laststarttime = ""
var controlParams map[string]interface{}

var fullDirList = []string{}
var SuggestionMap = []map[string]interface{}{}
var SuggestionList = []string{}
var CurrentMap = []map[string]interface{}{}
var suggestApp = []string{}
var filterMap = map[string][]string{"_app_": {"loki-ser"}}
var ipidcMap = map[string]string{}
var isHasKeywordGrep bool
var ipMapGlobal = map[string][]string{}
var logtypeMapGlobal = map[string][]string{}
var podIpMapGlobal = map[string][]string{}
var appenvMap = map[string]string{}
var fileMapping = map[string]string{}
var labelMappingGlobal = map[string][]string{}

var grepContextC = map[string]int{}
var grepContextB = map[string]int{}
var grepContextA = map[string]int{}
var grepContextD = map[string]int{}

type GrepContext struct {
	IsGrepContext bool   `json:"isGrepContext"`
	ContextValue  string `json:"contextValue"`
	ContextNum    int    `json:"contextNum"`
	Keyword       string `json:"keyword"`
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("tree"),
	readline.PcItem("login"),
	readline.PcItem("go",
		makeShortCut(user)...,
	),
	readline.PcItem("mode",
		readline.PcItem("root"),
		readline.PcItem("user"),
	),
	readline.PcItem("ls",
		readline.PcItem("-l"),
		readline.PcItem("-a"),
		readline.PcItem("-s"),
		readline.PcItem("-p"),
	),
	readline.PcItem("cd",
		makeCompleter()...,
	),
	readline.PcItem("more"),
	readline.PcItem("less"),
	readline.PcItem("tail"),
	readline.PcItem("history"),
	readline.PcItem("bye"),
	readline.PcItem("exit"),
	readline.PcItem("help"),
)

var moreCompleter = readline.PcItem("more",
	makeMoreCompletion(app_name)...,
)
var lsCompleter = readline.PcItem("ls",
	makeLsCompletion(app_name)...,
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

func makeShortCut(user string) []readline.PrefixCompleterInterface {
	args := []readline.PrefixCompleterInterface{}
	shortCutData := kernel.FetchData("application", user)
	for _, v := range shortCutData {
		rl := readline.PcItem(v)
		args = append(args, rl)
	}

	return args
}

func makeLsCompletion(app string) []readline.PrefixCompleterInterface {
	args := []readline.PrefixCompleterInterface{}
	labelMapping := map[string][]string{}
	labelMappingGlobal = map[string][]string{}

	env := appenvMap[user+":"+app]
	respData := kernel.GetLabelSeriesByApp(app, "", "", "", token, env)
	if len(respData) > 0 {
		for _, data := range respData {
			labelMap := data.(map[string]interface{})
			for label, value := range labelMap {
				val := value.(string)
				labelMapping[label] = append(labelMapping[label], val)
				labelMappingGlobal[user+":"+app+":"+label] = append(labelMappingGlobal[user+":"+app+":"+label], val)
			}
		}
	}

	for label, _ := range labelMapping {
		rl := readline.PcItem(label)
		args = append(args, rl)
	}

	return args
}

func makeCompleter() []readline.PrefixCompleterInterface {
	args := []readline.PrefixCompleterInterface{}

	applist := []string{}
	productMap := kernel.GetNode("root", user)
	fullDirList = kernel.FetchFullData()
	CurrentMap = productMap

	for index, _ := range productMap {
		applist = append(applist, productMap[index]["app"].(string))
	}
	SuggestionList = applist
	applist = append(applist, fullDirList...)

	for _, v := range applist {
		rl := readline.PcItem(v)
		args = append(args, rl)
	}
	return args
}

func makeMoreCompletion(app string) []readline.PrefixCompleterInterface {
	args := []readline.PrefixCompleterInterface{}
	ipMap := map[string][]string{}
	logtypeMap := map[string][]string{}
	podIpMap := map[string][]string{}
	ipMapGlobal = map[string][]string{}
	podIpMapGlobal = map[string][]string{}
	logtypeMapGlobal = map[string][]string{}
	// fmt.Println(app)
	env := appenvMap[user+":"+app]
	// fmt.Println(token)
	respData := kernel.GetLabelSeriesByApp(app, "", "", "", token, env)
	if len(respData) > 0 {
		for _, value := range respData {
			labelMap := value.(map[string]interface{})
			ip, ok_ip := labelMap["_ip_"].(string)
			filename, ok_file := labelMap["filename"].(string)
			podIp, ok_popip := labelMap["_pod_ip_"].(string)

			if !ok_file {
				promtMsg := fmt.Sprintf(warningCode, "文件路径对应的标签不为filename")
				fmt.Println(promtMsg)
			}
			if ok_ip {
				logtype, fileRealname := resolveLogTypeLabels(filename)
				logtypeMap[logtype] = append(logtypeMap[logtype], fileRealname)
				logtypeMapGlobal[logtype] = append(logtypeMapGlobal[logtype], filename)
				fileMapping[fileRealname] = filename

				ipMap[ip] = append(ipMap[ip], fileRealname)
				ipMapGlobal[ip] = append(ipMapGlobal[ip], filename)
			}
			if ok_popip {
				podIpMap[podIp] = append(podIpMap[podIp], "")
				podIpMapGlobal[podIp] = append(podIpMapGlobal[podIp], filename)
			}

		}
	}

	if len(ipMap) > 0 {
		args = filenameLinkage(ipMap, "ip", args)
		args = filenameLinkage(logtypeMap, "file", args)
	} else if len(podIpMap) > 0 {
		// fmt.Println(podIpMap)
		args = filenameLinkage(podIpMap, "pod", args)
	}

	return args
}

func filenameLinkage(filemap map[string][]string, flag string, args []readline.PrefixCompleterInterface) []readline.PrefixCompleterInterface {
	for i, v := range filemap {
		if flag == "pod" {
			rl := readline.PcItem(i)
			args = append(args, rl)
		} else {
			rl := readline.PcItem(i,
				makeFileCompletion(v, flag)...,
			)
			args = append(args, rl)
		}
	}
	return args
}

func makeFileCompletion(filelist []string, flag string) []readline.PrefixCompleterInterface {
	args := []readline.PrefixCompleterInterface{}
	// fmt.Println(filelist)
	filetemp := filelist
	if flag == "file" || flag == "ip" {
		filetemp = RemoveRepByMap(filelist)
	}

	if len(filetemp) > 0 {
		if !sort.StringsAreSorted(filetemp) {
			sort.Strings(filetemp)
		}
	}
	for _, v := range filetemp {
		rl := readline.PcItem(v)
		args = append(args, rl)
	}
	return args
}

func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func resolveLogTypeLabels(filename string) (string, string) {
	startSign := strings.LastIndex(filename, "/")
	fileRealname := filename[startSign+1:]

	endSignh := strings.Index(fileRealname, "-")
	endSignd := strings.Index(fileRealname, ".")
	endSign := 0
	if endSignh == -1 || endSignh > endSignd {
		endSign = endSignd
	} else {
		endSign = endSignh
	}
	logtype := fileRealname[:endSign]

	return logtype, fileRealname
}

func doFilterParamProcessing(args []string, orderIndex int) (filterParam map[string][]string, moreBehindArgsNum int) {
	filterParam = map[string][]string{}
	filterParam["_app_"] = append(filterParam["_app_"], app_name)

	for index, value := range args {
		if index == 0 {
			continue
		}
		moreBehindArgsNum = index
		if index != orderIndex {
			if index == 1 {
				ipOrLogtype := value
				isIp := checkIp(ipOrLogtype)
				if isIp {
					if len(ipMapGlobal) > 0 {
						filterParam["_ip_"] = append(filterParam["_ip_"], ipOrLogtype)
					} else {
						filterParam["_pod_ip_"] = append(filterParam["_pod_ip_"], ipOrLogtype)
					}
				} else {
					filelist := logtypeMapGlobal[ipOrLogtype]
					filetemp := RemoveRepByMap(filelist)
					filterParam["filename"] = append(filterParam["filename"], filetemp...)
				}
			}
			if index == 2 {
				fileRealname := args[2]
				filename := fileMapping[fileRealname]
				filterParam["filename"] = []string{filename}
			}
		}
	}
	// fmt.Println(filterParam)
	return
}

func checkArgsExistedOrderRule(args []string) (bool, int) {
	isPositive := false
	orderIndex := -1
	for index, value := range args {
		if value == "-asc" || value == "-desc" {
			orderIndex = index
			if value == "-asc" {
				isPositive = true
			}
			break
		}
	}
	return isPositive, orderIndex
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func grepContextStore(grepContext map[string]int, grepArgs []string, index int) {
	if len(grepArgs) > index+2 {
		cnum, err := strconv.Atoi(grepArgs[index+1])
		if err != nil {
			help(false)
		}
		grepContext[grepArgs[index+2]] = cnum
	} else {
		help(false)
	}
}

func Run(username string, app string) {

	//print welcome
	banner()

	//init args
	// token = kernel.FetchTokenByUsername(user)

	//user = username
	if username != "" {
		user = username
	}
	token = kernel.FetchTokenByUsername(user)

	if app != "" {
		app_name = app
		filterMap["_app_"] = []string{app_name}
		s_path = kernel.FetchPathByApp(app_name)
		nodeType = "application"
	}

	completer = readline.NewPrefixCompleter(
		readline.PcItem("go",
			makeShortCut(user)...,
		),
		readline.PcItem("cd",
			makeCompleter()...,
		),
		moreCompleter,
		lsCompleter,
	)

	prefix = "[" + user + "@lokishell " + s_path + "]$"

	l, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[32m" + prefix + "\033[0m ",
		HistoryFile:         "/tmp/readline.tmp",
		AutoComplete:        completer,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})

	if err != nil {
		panic(err)
	}
	defer l.Close()

	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	log.SetOutput(l.Stderr())

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				continue
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		isPositiveSequence := false
		isHasKeywordGrep = false
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				fmt.Println("vi")
				l.SetVimMode(true)
			case "emacs":
				fmt.Println("emacs")
				l.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case strings.HasPrefix(line, "cd "):
			o_path := s_path
			t_path = line[3:]
			cd(l, t_path)
			clearFilterMap(o_path, s_path)
		case strings.HasPrefix(line, "go "):
			o_path := s_path
			params := line[3:]
			paramlist := strings.Fields(params)
			appname := app_name
			env := "prod"
			for index, param := range paramlist {
				if index == 0 {
					appname = param
				}
				if index == 1 {
					env = param[1:]
				}
			}
			// fmt.Println(appname + ":" + env)
			appenvMap[user+":"+appname] = env
			if env != "prod" && env != "test" {
				promtMsg := fmt.Sprintf(warningCode, "环境参数输入不正确, 无法识别")
				fmt.Println(promtMsg)
				break
			}
			link(l, appname)
			clearFilterMap(o_path, s_path)
		case strings.HasPrefix(line, "filter "):
			args := strings.Fields(line)
			a := util.IsContain(args, "-a")
			d := util.IsContain(args, "-d")
			if !a && !d {
				filterdoc()
			} else {
				aIndex := util.IsElementExist(args, "-a")
				dIndex := util.IsElementExist(args, "-d")
				switch {
				case len(dIndex) > 0 && len(aIndex) > 0:
					aKeyword := args[aIndex[0]:]
					dKeyword := args[dIndex[0]:]
					if len(aKeyword) > 1 && len(dKeyword) > 1 {
						switch {
						case dIndex[0] > aIndex[0]+1:
							aFilters := args[aIndex[0]+1 : dIndex[0]]
							filterAdd(aFilters)
							dFilters := args[dIndex[0]+1:]
							filterDelete(dFilters)
							showFilter()
						case aIndex[0] > dIndex[0]+1:
							aFilters := args[aIndex[0]+1:]
							filterAdd(aFilters)
							dFilters := args[dIndex[0]+1 : aIndex[0]]
							filterDelete(dFilters)
							showFilter()
						default:
							filterdoc()
						}
					} else {
						filterdoc()
					}
				case len(aIndex) > 0 && len(dIndex) == 0:
					aFilters := args[aIndex[0]+1:]
					if len(aFilters) > 0 {
						filterAdd(aFilters)
						showFilter()
					} else {
						filterdoc()
					}
				case len(dIndex) > 0 && len(aIndex) == 0:
					dFilters := args[dIndex[0]+1:]
					if len(dFilters) > 0 {
						filterDelete(dFilters)
						showFilter()
					} else {
						filterdoc()
					}
				default:
					filterdoc()
				}
			}
		case strings.HasPrefix(line, "more "):
			var filterParam map[string][]string

			args := strings.Fields(line)
			//detect if contian keyword filter
			grep := util.IsContain(args, "grep")
			pipe := util.IsContain(args, "|")
			timerange := false

			orderParamIndex := -1
			moreBehindArgsNum := 0
			isHasKeywordGrep = true
			if !pipe && len(args) > 1 {
				// fmt.Println(args)
				isPositiveSequence, orderParamIndex = checkArgsExistedOrderRule(args)
				filterParam, _ = doFilterParamProcessing(args, orderParamIndex)

				isRelativeTime := true
				endtime = "1h"
				more(l, filterParam, starttime, endtime, "5", "20", "nokey", isPositiveSequence, false, false, isRelativeTime)
			} else {

				pipeIndex := util.IsElementExist(args, "|")
				firstIndex := pipeIndex[0]
				firstArgs := args[:firstIndex]
				// fmt.Println(firstArgs)
				isPositiveSequence, orderParamIndex = checkArgsExistedOrderRule(firstArgs)
				if len(firstArgs) > 1 {
					filterParam, moreBehindArgsNum = doFilterParamProcessing(firstArgs, orderParamIndex)
				} else {
					filterParam = filterMap
				}

				lastIndex := pipeIndex[len(pipeIndex)-1] + 1
				lastArgs := args[lastIndex:]
				if util.IsContain(lastArgs, "grep") {
					timerange = false
				} else {
					timerange = true
				}

				keywords := ""
				timepoint := false
				isRelativeTime := false
				if timerange && grep {
					//有时间过滤和关键字过滤
					grepArgs := args[2+moreBehindArgsNum : lastIndex-1]
					// fmt.Println(grepArgs)
					timeArgs := args[lastIndex:]
					// fmt.Println(timeArgs)
					timepoint = false

					if len(timeArgs) != 0 {
						if len(timeArgs) > 1 {
							//timerange
							starttime = timeArgs[0]
							endtime = timeArgs[1]
							// 判断endtime是不是时间点
							if strings.HasSuffix(endtime, "min") {
								timepoint = true
							} else {
								timepoint = false
							}
						} else {
							starttime = timeArgs[0]
							endtime = timeArgs[0]
							timepoint = true

							if strings.HasSuffix(endtime, "m") || strings.HasSuffix(endtime, "h") || strings.HasSuffix(endtime, "d") {
								isRelativeTime = true
								timepoint = false
							}
						}
					}

					grepContextC = map[string]int{}
					grepContextB = map[string]int{}
					grepContextA = map[string]int{}
					grepContextD = map[string]int{}

					isGrepParam := make([]bool, len(grepArgs))
					for index, value := range grepArgs {
						if index > 0 && isGrepParam[index-1] {
							isGrepParam[index] = false
							continue
						}
						switch value {
						case "-C":
							grepContextStore(grepContextC, grepArgs, index)
							isGrepParam[index] = true
							continue
						case "-B":
							grepContextStore(grepContextB, grepArgs, index)
							isGrepParam[index] = true
							continue
						case "-A":
							grepContextStore(grepContextA, grepArgs, index)
							isGrepParam[index] = true
							continue
						default:
							if value != "|" && value != "grep" && value != "" {
								grepContextD[value] = 0
							}
							isGrepParam[index] = false
						}

						if value != "|" && value != "grep" && value != "" {
							keywords += value + ","
						}
					}

					if len(keywords) > 0 {
						keywords = keywords[:len(keywords)-1]
					} else {
						keywords = "nokey"
					}
					// fmt.Println(keywords)

					isGrepContext := false
					for _, value := range isGrepParam {
						if value {
							isGrepContext = value
							break
						}
					}

					more(l, filterParam, starttime, endtime, "5", "20", keywords, isPositiveSequence, timepoint, isGrepContext, isRelativeTime)
				} else if timerange && grep == false {
					//有时间过滤,但是没关键字过滤
					timeArgs := args[lastIndex:]
					// fmt.Println(timeArgs)
					timepoint = false
					if len(timeArgs) != 0 && len(timeArgs) <= 2 {
						if len(timeArgs) > 1 {
							//timerange
							starttime = timeArgs[0]
							endtime = timeArgs[1]
							if strings.HasSuffix(endtime, "min") {
								timepoint = true
							} else {
								timepoint = false
							}
						} else {
							starttime = timeArgs[0]
							endtime = timeArgs[0]
							timepoint = true

							if strings.HasSuffix(endtime, "m") || strings.HasSuffix(endtime, "h") || strings.HasSuffix(endtime, "d") {
								isRelativeTime = true
								timepoint = false
							}
						}

						more(l, filterParam, starttime, endtime, "5", "20", "nokey", isPositiveSequence, timepoint, false, isRelativeTime)
					} else {
						help(false)
					}
				} else if timerange == false && grep {
					grepContextC = map[string]int{}
					grepContextB = map[string]int{}
					grepContextA = map[string]int{}
					grepContextD = map[string]int{}

					//没有时间过滤，但是有关键字过滤,默认当前时间之前的1小时
					grepArgs := args[2+moreBehindArgsNum:]
					// fmt.Println(args)
					// fmt.Println(grepArgs)
					keywords := ""
					isGrepParam := make([]bool, len(grepArgs))
					for index, value := range grepArgs {
						if index > 0 && isGrepParam[index-1] {
							isGrepParam[index] = false
							continue
						}
						switch value {
						case "-C":
							grepContextStore(grepContextC, grepArgs, index)
							isGrepParam[index] = true
							continue
						case "-B":
							grepContextStore(grepContextB, grepArgs, index)
							isGrepParam[index] = true
							continue
						case "-A":
							grepContextStore(grepContextA, grepArgs, index)
							isGrepParam[index] = true
							continue
						default:
							if value != "|" && value != "grep" && value != "" {
								grepContextD[value] = 0
							}
							isGrepParam[index] = false
						}

						if value != "|" && value != "grep" && value != "" {
							keywords += value + ","
						}
					}

					if len(keywords) > 0 {
						keywords = keywords[:len(keywords)-1]
					} else {
						keywords = "nokey"
					}
					// fmt.Println(keywords)

					isGrepContext := false
					for _, value := range isGrepParam {
						if value {
							isGrepContext = value
							break
						}
					}

					// dt := time.Now()
					// endtime = dt.Format(layout)
					// d, _ := time.ParseDuration("-24h")
					// starttime = dt.Add(d).Format(layout)
					isRelativeTime = true
					endtime = "1h"
					more(l, filterParam, starttime, endtime, "5", "20", keywords, isPositiveSequence, false, isGrepContext, isRelativeTime)
				}

			}

		case strings.HasPrefix(line, "seek "):
			// [seek] 100 | traceid[filename/ip] | c 100
			args := strings.Split(line, "|")
			id := ""
			mode := "filename"
			traceid := ""
			border := "c"
			upnum := "30"
			downnum := "5"
			for index, param := range args {
				element := strings.TrimSpace(param)
				content := strings.Fields(element)
				key := content[0]
				if key == "seek" && len(content) > 1 {
					id = content[1]
				}
				if index == 1 {
					_, err := strconv.Atoi(key)
					if err != nil {
						switch key {
						case "filename":
							mode = "filename"
						case "ip":
							mode = "ip"
						case "applicaton":
							mode = "applicaton"
						default:
							mode = "traceid"
							traceid = key
						}
					} else {
						upnum = key
					}
				}
				if index == 2 && len(content) > 1 {
					border = key
					upnum = content[1]
				} else if index == 2 && len(content) == 1 {
					_, err := strconv.Atoi(key)
					if err != nil {
						border = key
					} else {
						upnum = key
					}
				}
			}
			seek(id, mode, traceid, border, upnum, downnum)

		case strings.HasPrefix(line, "less "):

		case strings.HasPrefix(line, "tail "):

			pipe := strings.Contains(line, "|")
			if pipe {
				grep := strings.Contains(line, "grep")
				grepIndex := strings.Index(line, "grep") + 4

				keywordStr := line[grepIndex:]
				keywordStr = strings.TrimSpace(keywordStr)
				keyList := strings.Split(keywordStr, " ")
				keywords := ""
				for _, value := range keyList {
					keywords += value + ","
				}
				keywords = keywords[:len(keywords)-1]

				if grep && keywords != "" && keywords != "," {
					tail(filterMap, true, keywords, "prod")
				} else {
					tail(filterMap, true, "nokey", "prod")
				}
			} else {
				help(false)
			}

		case strings.HasPrefix(line, "ls "):
			arrCommandStr := strings.Fields(line)
			labels := []string{"_ip_", "_pod_ip_"}
			if len(arrCommandStr) > 1 {
				labels = []string{arrCommandStr[1]}
			}
			listMode := false
			listFilename := false
			if util.IsContain(arrCommandStr, "-l") {
				listMode = true
				ls(listMode, listFilename, labels)
			} else if util.IsContain(arrCommandStr, "-f") {
				listFilename = true
				ls(false, listFilename, labels)
			} else {
				listMode = false
				ls(listMode, listFilename, labels)
			}
		case line == "mode":
			if l.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case line == "ls":
			labels := []string{"_ip_", "_pod_ip_"}
			ls(false, false, labels)
		case strings.HasPrefix(line, "back"):
			displayRecentOnceQueryResult()
		case line == "filter":
			filterdoc()
		case line == "more":
			// filelist := logtypeMapGlobal["application"]
			// filterMap["filename"] = append(filterMap["filename"], filelist...)
			isHasKeywordGrep = true

			// dt := time.Now()
			// endtime = dt.Format(layout)
			// d, _ := time.ParseDuration("-1h")
			// starttime = dt.Add(d).Format(layout)
			isRelativeTime := true
			endtime = "1h"
			more(l, filterMap, starttime, endtime, "5", "20", "nokey", isPositiveSequence, false, false, isRelativeTime)
		case line == "less":
			dt := time.Now()
			endtime = dt.Format(layout)
			d, _ := time.ParseDuration("-1h")
			starttime = dt.Add(d).Format(layout)
			less(l, filterMap, starttime, endtime, "5", "500", "nokey", false, false)
		case line == "help":
			help(true)

		case line == "history":
			f, err := ioutil.ReadFile("/tmp/readline.tmp")
			if err != nil {
				fmt.Println("read fail", err)
			}
			fmt.Printf(string(f))
		case line == "tail":
			tail(filterMap, false, "nokey", "devops")
		case line == "login":
			login(l, "nouser")
		case line == "tree":
			usage(l.Stderr())
		case strings.HasPrefix(line, "say"):
			line := strings.TrimSpace(line[3:])
			if len(line) == 0 {
				log.Println("say what?")
				break
			}
			// go func() {
			// 	for range time.Tick(time.Second) {
			// 		log.Println(line)
			// 	}
			// }()
		case line == "pwd":
			fmt.Println(s_path)
		case line == "bye":
			goto exit
		case line == "exit":
			goto exit
		case line == "quit":
			goto exit
		case line == "sleep":
			log.Println("sleep 4 second")
			time.Sleep(4 * time.Second)
		case line == "":
		default:
			// id upnum downnum | traceid
			argsStr := strings.TrimSpace(line)
			args := strings.Split(argsStr, "|")

			uid := ""
			traceid := ""
			upnum := "30"
			downnum := "5"
			if len(args) > 0 {
				idargs := strings.TrimSpace(args[0])
				paramlist := strings.Fields(idargs)
				for index, value := range paramlist {
					if index == 0 {
						uid = value
					}
					if index == 1 {
						upnum = value
					}
					if index == 2 {
						downnum = value
					}
				}
			}
			if len(args) == 2 {
				traceargs := strings.TrimSpace(args[1])
				if traceargs != "" {
					traceid = traceargs
				}
			}

			// fmt.Println(uid + ":" + linenum)
			seek(uid, "filename", traceid, "c", upnum, downnum)
		}
	}
exit:
}
