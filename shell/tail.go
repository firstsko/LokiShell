package shell

import (
	"context"
	"fmt"
	"log"
	"lokishell/kernel"
	"lokishell/util"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sacOO7/gowebsocket"
)

const ali_server = "loki-ali.hfinside.com:8534"
const up_server = "loki.hfinside.com:8502"
const test_server = "10.64.0.213:80"

var appQuery string
var keyQuery string
var otherQuery string

type Stream struct {
	Streams []struct {
		Stream struct {
			Filename string `json:"filename"`
			App      string `json:"_app_"`
			Appid    string `json:"_appid_"`
		} `json:"stream"`
		Values [][]string `json:"values"`
	} `json:"streams"`
}

func stringBuilder(server string, ipQuery string) string {
	var url strings.Builder
	url.WriteString("ws://")
	url.WriteString(server)
	url.WriteString("/loki/api/v1/tail?query=")
	url.WriteString("{")
	url.WriteString(appQuery)
	url.WriteString(ipQuery)
	url.WriteString(otherQuery)
	url.WriteString("}")
	url.WriteString(keyQuery)
	return url.String()
}

func tail(filterMap map[string][]string, grep bool, keyword string, tenant string) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	aliIpList := []string{}
	upIpList := []string{}

	app := "loki-ser"
	if appList, ok := filterMap["app"]; ok {
		app = appList[0]
	}
	if appList, ok := filterMap["_app_"]; ok {
		app = appList[0]
	}
	env := appenvMap[user+":"+app]
	// fmt.Println(env)

	if env != "test" {
		instanceMap := kernel.GetInstace(app)
		for key, valueList := range filterMap {
			if strings.Contains(key, "ip") {
				if len(valueList) > 0 {
					for _, ip := range valueList {
						ipFlag := false
						for index, _ := range instanceMap {
							ipdb := instanceMap[index]["internal_ip"].(string)
							if ip == ipdb {
								idc := instanceMap[index]["idc"].(string)
								switch idc {
								case "aliyun":
									aliIpList = append(aliIpList, ip)
								case "upyun":
									upIpList = append(upIpList, ip)
								default:
									upIpList = append(upIpList, ip)
								}
								ipFlag = true
							}
						}
						if !ipFlag {
							fmt.Println("current app not exist ip ", ip)
						}
					}
				}
			}
		}
	}

	// if len(aliIpList) == 0 && len(upIpList) == 0 {
	// 	for index, _ := range instanceMap {
	// 		ip := instanceMap[index]["internal_ip"].(string)
	// 		idc := instanceMap[index]["idc"].(string)
	// 		switch idc {
	// 		case "aliyun":
	// 			aliIpList = append(aliIpList, ip)
	// 		case "upyun":
	// 			upIpList = append(upIpList, ip)
	// 		default:
	// 			upIpList = append(upIpList, ip)
	// 		}
	// 	}
	// }

	appQuery = "_app_=\"" + app + "\""

	aliIpQuery := strings.Builder{}
	upIpQuery := strings.Builder{}
	if env != "test" {
		if len(aliIpList) > 0 {
			aliIpQuery.WriteString(",_ip_=~\"")
			for index, ip := range aliIpList {
				if index == 0 {
					aliIpQuery.WriteString(ip)
				} else {
					aliIpQuery.WriteString("|")
					aliIpQuery.WriteString(ip)
				}
			}
			aliIpQuery.WriteString("\"")
		}
		if len(upIpList) > 0 {
			upIpQuery.WriteString(",_ip_=~\"")
			for index, ip := range upIpList {
				if index == 0 {
					upIpQuery.WriteString(ip)
				} else {
					upIpQuery.WriteString("|")
					upIpQuery.WriteString(ip)
				}
			}
			upIpQuery.WriteString("\"")
		}
	}

	//默认查询application类型的日志
	filelist := logtypeMapGlobal["application"]
	filetemp := RemoveRepByMap(filelist)
	filterMap["filename"] = append(filterMap["filename"], filetemp...)

	otherQuery = ""
	for key, valueList := range filterMap {
		if key == "app" || key == "_app_" || (env != "test" && (key == "ip" || key == "_ip_")) {
			continue
		}
		if len(valueList) > 0 {
			otherQuery += "," + key + "=~\""
			for index, value := range valueList {
				if index == 0 {
					otherQuery += value
				} else {
					otherQuery += "|" + value
				}
			}
			otherQuery += "\""
		}
	}

	keyQuery = ""
	if keyword != "nokey" {
		keywords := strings.Split(keyword, ",")
		keyQuery += "|~\""
		for index, key := range keywords {
			if index == 0 {
				keyQuery += key
			} else {
				keyQuery += "|" + key
			}
		}
		keyQuery += "\""
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	tenant = kernel.GetTenantIdByAppName(app_name, token)
	// fmt.Println(tenant)

	if env == "test" {
		test_url := stringBuilder(test_server, "")

		socket := gowebsocket.New(test_url)
		createWebsocketExchange(socket, keyword, "test", interrupt, ctx, cancel, tenant)
	} else {
		ali_url := stringBuilder(ali_server, aliIpQuery.String())
		up_url := stringBuilder(up_server, upIpQuery.String())
		// fmt.Println(ali_url)
		// fmt.Println(up_url)

		if (len(aliIpList) == 0 && len(upIpList) == 0) || (len(aliIpList) > 0 && len(upIpList) > 0) {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				socket := gowebsocket.New(ali_url)
				createWebsocketExchange(socket, keyword, "aliyun", interrupt, ctx, cancel, tenant)
			}()
			go func() {
				defer wg.Done()
				socket := gowebsocket.New(up_url)
				createWebsocketExchange(socket, keyword, "upyun", interrupt, ctx, cancel, tenant)
			}()
			wg.Wait()
		} else if aliIpQuery.String() != "" {
			socket := gowebsocket.New(ali_url)
			createWebsocketExchange(socket, keyword, "aliyun", interrupt, ctx, cancel, tenant)
		} else {
			socket := gowebsocket.New(up_url)
			createWebsocketExchange(socket, keyword, "upyun", interrupt, ctx, cancel, tenant)
		}
	}

}

func createWebsocketExchange(socket gowebsocket.Socket, keyword string, category string, interrupt chan os.Signal, ctx context.Context, cancel context.CancelFunc, tenant string) {
	socket.RequestHeader.Set("X-Scope-OrgID", tenant)

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server ", category)
	}

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", category, err)
	}

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		//fmt.Printf("%T\n", message)
		msg, _ := util.JsonToMap(message)

		for _, item := range msg["streams"].([]interface{}) {
			for _, event := range item.(map[string]interface{})["values"].([]interface{}) {

				var events = []string{}
				for _, row := range event.([]interface{}) {
					events = append(events, row.(string))

					if len(events) == 2 {
						timestamp := events[0]
						unixIntValue, _ := strconv.ParseInt(timestamp, 10, 64)
						t := time.Unix(0, unixIntValue)
						timeDisplay := t.Format("2006-01-02 15:04:05")
						timeDisplay = fmt.Sprintf(infoCode, timeDisplay)

						if keyword != "nokey" {
							output := events[1]
							//hightlight keyword
							keywords := strings.Split(keyword, ",")
							for _, key := range keywords {
								hightlightKey := fmt.Sprintf(warningCode, key)
								output = strings.Replace(output, key, hightlightKey, -1)
							}
							fmt.Println(timeDisplay + " " + output)
						} else {
							fmt.Println(timeDisplay + " " + events[1])
						}

					}

				}
			}
		}

		//fmt.Println(msgMap)
		//log.Println("Recieved message " + msgMap)
	}

	socket.OnBinaryMessage = func(data []byte, socket gowebsocket.Socket) {
		//log.Println("Recieved binary data ", data)
	}

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		//log.Println("Recieved ping " + data)
	}

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		//log.Println("Recieved pong " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		// log.Println("Disconnected from server ", category)
		return
	}
	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt server")
			socket.Close()
			cancel()
			return
		case <-ctx.Done():
			socket.Close()
			return
		}
	}
}
