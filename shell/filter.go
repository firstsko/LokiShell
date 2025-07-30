package shell

import (
	"fmt"
	"lokishell/kernel"
	"net"
	"os"
	"regexp"
	"strings"
)

func checkIp(ip string) bool {
	address := net.ParseIP(ip)
	return address != nil
}

func checkFilename(filename string) bool {
	fileReg := `.(log|xxx)$`
	re, _ := regexp.Compile(fileReg)
	match := re.MatchString(filename)
	return match
}

func filterAdd(filters []string) {
	// 1.判断filter是什么标签；_appid_、_ip_、_env_、_idc_、filename、_app_
	// {"_app_":"acp-server","_appid_":"acp.acp-server","_ip_":"172.19.252.166","_env_":"prod","_idc":"aliyun",
	// "filename":"/app/docker/acp-server/logs/100003171/apollo-configservice.log"}
	if len(filters) > 0 {

		for _, filter := range filters {
			tags := strings.Split(filter, "=")
			if len(tags) == 2 {
				tagValue := strings.Split(tags[1], ",")
				instanceMap := kernel.GetInstace(app_name)

				isEffectiveIp := false
				for _, ip := range tagValue {
					for index, _ := range instanceMap {
						ipdb := instanceMap[index]["internal_ip"].(string)
						if ip == ipdb {
							isEffectiveIp = true
						}
					}
				}
				if !isEffectiveIp {
					fmt.Println("当前指定的ip为无效ip,请确认是否为当前应用下的ip")
					return
				}

				filterMap[tags[0]] = append(filterMap[tags[0]], tagValue...)
			} else {
				filterdoc()

			}
		}
	}
}

func filterDelete(filters []string) {
	if len(filters) > 0 {
		for _, filter := range filters {
			tags := strings.Split(filter, "=")
			if len(tags) == 2 {
				if _, ok := filterMap[tags[0]]; ok {
					tagValue := strings.Split(tags[1], ",")
					for _, inputValue := range tagValue {
						for index, filterValue := range filterMap[tags[0]] {
							if inputValue == filterValue {
								filterMap[tags[0]] = append(filterMap[tags[0]][:index], filterMap[tags[0]][index+1:]...)
								if len(filterMap[tags[0]]) == 0 {
									delete(filterMap, tags[0])
								}
							}
						}
					}
				} else {
					checkmsg := fmt.Sprintf(warningCode, tags[0]+" 过滤器不存在，无法删除")
					fmt.Fprintln(os.Stdout, checkmsg)
				}

			} else {
				filterdoc()
			}
		}
	}

}

func clearFilterMap(o_path string, s_path string) {
	if o_path != s_path {
		for key := range filterMap {
			if key == "_app_" || key == "app" {
				continue
			}
			delete(filterMap, key)
		}
	}
}

func showFilter() {
	for key, value := range filterMap {
		fmt.Println(key, ":", value)
	}
}
