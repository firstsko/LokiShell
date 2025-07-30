package shell

import (
	"fmt"
	"lokishell/kernel"
	"sort"
)

func ls(listMode bool, listFilename bool, labels []string) {
	if nodeType != "application" {
		if listMode == true {
			for index, _ := range CurrentMap {
				fmt.Println(CurrentMap[index]["app"].(string) + "[" + CurrentMap[index]["alias"].(string) + "]")
			}
		} else {
			listStr := ""
			for index, _ := range CurrentMap {
				listStr += CurrentMap[index]["app"].(string) + "[" + CurrentMap[index]["alias"].(string) + "]  "
			}
			fmt.Println(listStr)
		}
	} else {
		if listFilename {
			ip := ""
			starttime := ""
			endtime := ""
			env := appenvMap[user+":"+app_name]
			respData := kernel.GetLabelSeriesByApp(app_name, ip, starttime, endtime, token, env)
			if len(respData) > 0 {
				for _, value := range respData {
					labelMap := value.(map[string]interface{})
					ip, ok_ip := labelMap["_ip_"].(string)
					filename, ok_file := labelMap["filename"].(string)
					podIp, ok_popip := labelMap["_pod_ip_"].(string)
					if ok_ip && ok_file {
						fmt.Println(ip + " : " + filename)
					}
					if ok_popip && ok_file {
						fmt.Println(podIp + " : " + filename)
					}
				}
			} else {
				promtMsg := fmt.Sprintf(warningCode, "当前查询无相关结果")
				fmt.Println(promtMsg)
			}

		} else {

			if len(labelMappingGlobal) > 0 {
				for _, la := range labels {
					labelvals, ok := labelMappingGlobal[user+":"+app_name+":"+la]
					if ok {
						vallist := RemoveRepByMap(labelvals)
						if len(vallist) > 0 {
							if !sort.StringsAreSorted(vallist) {
								sort.Strings(vallist)
							}
						}
						for _, val := range vallist {
							fmt.Println(val)
						}
					}
				}

			} else {
				promtMsg := fmt.Sprintf(warningCode, "当前查询无相关结果")
				fmt.Println(promtMsg)
			}

			// instanceMap := kernel.GetInstace(app_name)
			// for index, _ := range instanceMap {
			// 	fmt.Println(instanceMap[index]["internal_ip"].(string))
			// }
		}
	}

}
