package kernel

import (
	"context"
	"fmt"
	"lokishell/util"
	"strconv"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	// "bytes"
)

const url_config = "http://api.hoc.huifu.com"

// const url_config = "http://localhost:5100"

// const url_event = "http://api.loki.hfinside.com"
const redCode = "\033[31m %s \033[0m"
const warningCode = "\033[33m %s \033[0m"
const successCode = "\033[32m %s \033[0m"
const infoCode = "\033[36m %s \033[0m"
const logcli = "./logcli"

var product = []string{}
var system = []string{}
var application = []string{}
var fullDirList = []string{}
var wg sync.WaitGroup

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

type ServiceTreeGenerated struct {
	ReqID   string `json:"ReqId"`
	Project string `json:"Project"`
	Action  string `json:"Action"`
	RetCode int    `json:"RetCode"`
	Data    []Data `json:"Data"`
}

type PathGenerated struct {
	ReqID   string `json:"ReqId"`
	Project string `json:"Project"`
	Action  string `json:"Action"`
	RetCode int    `json:"RetCode"`
	Data    string `json:"Data"`
}

type InstanceGenerated struct {
	ReqID   string `json:"ReqId"`
	Project string `json:"Project"`
	Action  string `json:"Action"`
	RetCode int    `json:"RetCode"`
	Sets    []struct {
		InternalIP string `json:"internal_ip"`
		AppName    string `json:"app_name"`
		Port       string `json:"port"`
		TreeID     int    `json:"tree_id"`
		ArmsID     int    `json:"arms_id"`
		AppID      string `json:"app_id"`
		Idc        string `json:"idc"`
	} `json:"Sets"`
}

type Node struct {
	ID                   int         `json:"id"`
	AppID                string      `json:"app_id"`
	Name                 string      `json:"name"`
	Alias                string      `json:"alias"`
	Path                 string      `json:"path"`
	Parent               int         `json:"parent"`
	NodeType             int         `json:"node_type"`
	Children             []Node      `json:"children"`
	Title                string      `json:"title"`
	Expand               bool        `json:"expand"`
	Zone                 string      `json:"zone,omitempty"`
	ShortAppID           string      `json:"short_app_id,omitempty"`
	ProductID            int         `json:"product_id,omitempty"`
	Tags                 string      `json:"tags,omitempty"`
	GitHTTP              string      `json:"git_http,omitempty"`
	GitSSH               string      `json:"git_ssh,omitempty"`
	Git                  string      `json:"git,omitempty"`
	Svn                  string      `json:"svn,omitempty"`
	GitID                string      `json:"git_id,omitempty"`
	GitProjectID         int         `json:"git_project_id,omitempty"`
	RepoActive           int         `json:"repo_active,omitempty"`
	CreateTime           string      `json:"create_time,omitempty"`
	UpdateTime           time.Time   `json:"update_time,omitempty"`
	ReleaseStartTime     string      `json:"release_start_time,omitempty"`
	ReleaseEndTime       string      `json:"release_end_time,omitempty"`
	Branch               string      `json:"branch,omitempty"`
	InitID               int         `json:"init_id,omitempty"`
	Test                 int         `json:"test,omitempty"`
	TbProjectID          interface{} `json:"tb_project_id,omitempty"`
	TbTestID             interface{} `json:"tb_test_id,omitempty"`
	TbDirID              interface{} `json:"tb_dir_id,omitempty"`
	ManualExample        int         `json:"manual_example,omitempty"`
	AutoExample          interface{} `json:"auto_example,omitempty"`
	InterfaceAutoExample interface{} `json:"interface_auto_example,omitempty"`
	Example              interface{} `json:"example,omitempty"`
	Malfunction          interface{} `json:"malfunction,omitempty"`
	Team                 string      `json:"team,omitempty"`
	Dep                  string      `json:"dep,omitempty"`
	IsSonar              int         `json:"is_sonar,omitempty"`
	Mapping              int         `json:"mapping,omitempty"`
	DroneVersion         string      `json:"drone_version,omitempty"`
	IsHcm                int         `json:"is_hcm,omitempty"`
	DeployMode           string      `json:"deploy_mode,omitempty"`
	DefaultMaster        interface{} `json:"default_master,omitempty"`
	ArmsGreyInterval     interface{} `json:"arms_grey_interval,omitempty"`
	ArmsProiority        interface{} `json:"arms_proiority,omitempty"`
	IsGraydeploy         int         `json:"is_graydeploy,omitempty"`
	ArchType             interface{} `json:"arch_type,omitempty"`
	DingtalkToken        interface{} `json:"dingtalk_token,omitempty"`
	DepID                interface{} `json:"dep_id,omitempty"`
	ArmsDeployMode       interface{} `json:"arms_deploy_mode,omitempty"`
	CatSystemID          interface{} `json:"cat_system_id,omitempty"`
	IsTestReport         int         `json:"is_test_report,omitempty"`
	AutoPipeline         int         `json:"auto_pipeline,omitempty"`
	DefaultPipelineID    interface{} `json:"default_pipeline_id,omitempty"`
	HuifucloudProFlag    interface{} `json:"huifucloud_pro_flag,omitempty"`
	HuifuaopPro          interface{} `json:"huifuaop_pro,omitempty"`
	Prom                 int         `json:"prom,omitempty"`
	StartNotify          int         `json:"start_notify,omitempty"`
	StopNotify           int         `json:"stop_notify,omitempty"`
	SuccessNotify        int         `json:"success_notify,omitempty"`
	FailedNotify         int         `json:"failed_notify,omitempty"`
	HealthCheck          int         `json:"health_check,omitempty"`
	UniqueAppID          string      `json:"unique_app_id,omitempty"`
	Production           string      `json:"production,omitempty"`
	ArmsID               int         `json:"arms_id,omitempty"`
	DingtalkURL          interface{} `json:"dingtalk_url,omitempty"`
	IsDeleted            interface{} `json:"is_deleted,omitempty"`
	GreyInterval         interface{} `json:"grey_interval,omitempty"`
	Readness             interface{} `json:"readness,omitempty"`
	LivenessPath         interface{} `json:"liveness_path,omitempty"`
	Middleware           interface{} `json:"middleware,omitempty"`
	Language             interface{} `json:"language,omitempty"`
	JenkinsBuildTask     interface{} `json:"jenkins_build_task,omitempty"`
	DeployAddress        interface{} `json:"deploy_address,omitempty"`
	JacocoPort           interface{} `json:"jacoco_port,omitempty"`
	HTTPPort             interface{} `json:"http_port,omitempty"`
	TargetPath           interface{} `json:"target_path,omitempty"`
	BuildCmd             interface{} `json:"build_cmd,omitempty"`
	BuildType            interface{} `json:"build_type,omitempty"`
	StoreType            interface{} `json:"store_type,omitempty"`
	JdkVersion           interface{} `json:"jdk_version,omitempty"`
	IsCat                interface{} `json:"is_cat,omitempty"`
	IsLens               interface{} `json:"is_lens,omitempty"`
	ProductNameEn        interface{} `json:"product_name_en,omitempty"`
	ProductName          interface{} `json:"product_name,omitempty"`
	Status               interface{} `json:"status,omitempty"`
	DevOwner             interface{} `json:"dev_owner,omitempty"`
	DevOwnerCn           interface{} `json:"dev_owner_cn,omitempty"`
	TbMount              interface{} `json:"tb_mount,omitempty"`
	Version              int         `json:"version,omitempty"`
}

type Data struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	AppID    string `json:"app_id"`
	Path     string `json:"path"`
	NodeType int    `json:"node_type"`
	Parent   int    `json:"parent"`
	Children []Node `json:"children"`
}

type QueryRangeGenerate struct {
	ReqID   string `json:"ReqId"`
	Project string `json:"Project"`
	Action  string `json:"Action"`
	Data    struct {
		Status  string `json:"status"`
		RetCode int    `json:retCode`
		Message string `json.message`
		Data    struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Stream map[string]string `json:"stream"`
				Values [][]interface{}   `json:"values"`
			} `json:"result"`
			Stats struct {
				Summary struct {
					BytesProcessedPerSecond int     `json:"bytesProcessedPerSecond"`
					LinesProcessedPerSecond int     `json:"linesProcessedPerSecond"`
					TotalBytesProcessed     int     `json:"totalBytesProcessed"`
					TotalLinesProcessed     int     `json:"totalLinesProcessed"`
					ExecTime                float64 `json:"execTime"`
				} `json:"summary"`
				Store struct {
					TotalChunksRef        int `json:"totalChunksRef"`
					TotalChunksDownloaded int `json:"totalChunksDownloaded"`
					ChunksDownloadTime    int `json:"chunksDownloadTime"`
					HeadChunkBytes        int `json:"headChunkBytes"`
					HeadChunkLines        int `json:"headChunkLines"`
					DecompressedBytes     int `json:"decompressedBytes"`
					DecompressedLines     int `json:"decompressedLines"`
					CompressedBytes       int `json:"compressedBytes"`
					TotalDuplicates       int `json:"totalDuplicates"`
				} `json:"store"`
				Ingester struct {
					TotalReached       int `json:"totalReached"`
					TotalChunksMatched int `json:"totalChunksMatched"`
					TotalBatches       int `json:"totalBatches"`
					TotalLinesSent     int `json:"totalLinesSent"`
					HeadChunkBytes     int `json:"headChunkBytes"`
					HeadChunkLines     int `json:"headChunkLines"`
					DecompressedBytes  int `json:"decompressedBytes"`
					DecompressedLines  int `json:"decompressedLines"`
					CompressedBytes    int `json:"compressedBytes"`
					TotalDuplicates    int `json:"totalDuplicates"`
				} `json:"ingester"`
			} `json:"stats"`
		} `json:"data"`
	} `json:"Data"`
}

type LabelsGenerated struct {
	ReqID   string `json:"ReqId"`
	Project string `json:"Project"`
	Action  string `json:"Action"`
	Data    struct {
		Status string `json:"status"`
		Data   []struct {
			Idc       string `json:"_idc_"`
			IP        string `json:"_ip_"`
			Filename  string `json:"filename"`
			App       string `json:"_app_"`
			Appid     string `json:"_appid_"`
			Env       string `json:"_env_"`
			Tag       string `json:"_tag_"`
			Name      string `json:"__name__"`
			PodHostIP string `json:"_pod_host_ip_"`
			PodIP     string `json:"_pod_ip_"`
			PodName   string `json:"_pod_name_"`
			PodUID    string `json:"_pod_uid_"`
		} `json:"data"`
	} `json:"Data"`
}

type GatewayTokeGenerated struct {
	ReqID   string `json:"ReqId"`
	Project string `json:"Project"`
	Action  string `json:"Action"`
	RetCode int    `json:"RetCode"`
	Data    struct {
		UserID        int    `json:"user_id"`
		Avatar        string `json:"avatar"`
		DepParentName string `json:"dep_parent_name"`
		DepName       string `json:"dep_name"`
		UserEmail     string `json:"user_email"`
		Verify        string `json:"verify"`
		Inited        int    `json:"inited"`
		Admin         int    `json:"admin"`
	} `json:"Data"`
	UserName string `json:"UserName"`
	Token    string `json:"Token"`
}

func SetColor(msg string, conf, bg, text int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, text, msg, 0x1B)
}

//根据用户名获取网关token

func FetchTokenByUsername(username string) string {
	var response GatewayTokeGenerated

	request := gorequest.New().Timeout(10000 * time.Millisecond)
	request.Post(url_config).
		Set("Content-Type", "application/json").
		Param("Action", "HuifuLdapLogin").
		Param("Project", "Generally").
		Param("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com").
		Param("Signature", "f655e89631c4173ec08438f5b8a70faa841e4a8c").
		Param("UserName", username).
		Param("PassWord", "").
		EndStruct(&response)

	return response.Token
}

//根据应用名获得绝对路径

func FetchPathByApp(app string) string {
	var response PathGenerated

	request := gorequest.New().Timeout(10000 * time.Millisecond)
	request.Get(url_config).
		Set("Content-Type", "application/json").
		Send(`{
	"Action":"GetPathByApp",
	"Project": "Tree",
	"App":"` + app + `",
	"PublicKey": "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com",
	"Signature": "f655e89631c4173ec08438f5b8a70faa841e4a8c"
}`).
		EndStruct(&response)

	return response.Data
}

//获取全部服务树路径列表
func FetchFullData() []string {
	var response ServiceTreeGenerated

	request := gorequest.New().Timeout(10000 * time.Millisecond)

	request.Get(url_config).
		Set("Content-Type", "application/json").
		Send(`{
		"Action":"GetTree",
		"Project": "Tree",
		"PublicKey": "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com",
		"Signature": "f655e89631c4173ec08438f5b8a70faa841e4a8c"
	}`).
		EndStruct(&response)

	for index, _ := range response.Data {
		productElement := "/" + response.Data[index].Name
		fullDirList = append(fullDirList, productElement)
		for system, _ := range response.Data[index].Children {
			systemElement := "/" + response.Data[index].Name + "/" + response.Data[index].Children[system].Name
			fullDirList = append(fullDirList, systemElement)
			for app, _ := range response.Data[index].Children[system].Children {
				appElement := "/" + response.Data[index].Name + "/" + response.Data[index].Children[system].Name + "/" + response.Data[index].Children[system].Children[app].Name
				fullDirList = append(fullDirList, appElement)
			}
		}
	}
	return fullDirList
}

//获取当前目录下的列表
func FetchData(nodeType string, userName string) []string {
	var response ServiceTreeGenerated

	action := ""
	switch nodeType {
	case "product":
		action = "GetProductNode"
	case "system":
		action = "GetSystemNode"
	case "application":
		action = "GetApplicationNode"
	}

	request := gorequest.New().Timeout(10000 * time.Millisecond)
	request.Get(url_config).
		Set("Content-Type", "application/json").
		Send(`{
			"Action":"` + action + `",
			"Project": "Tree",
			"UserName":"` + userName + `",
			"PublicKey": "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com",
			"Signature": "f655e89631c4173ec08438f5b8a70faa841e4a8c"
	}`).
		EndStruct(&response)

	dirList := []string{}
	for index, _ := range response.Data {
		element := response.Data[index].Name
		dirList = append(dirList, element)
	}

	return dirList
}

//根据应用获取实例
func GetInstace(app string) []map[string]interface{} {

	var response InstanceGenerated
	request := gorequest.New().Timeout(10000 * time.Millisecond)

	request.Post(url_config).
		Set("Content-Type", "application/json").
		Param("Action", "GetInstanceListByAppName").
		Param("Project", "Tree").
		Param("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com").
		Param("Signature", "f655e89631c4173ec08438f5b8a70faa841e4a8c").
		Param("AppName", app).
		EndStruct(&response)

	var DataMap = []map[string]interface{}{}
	for index, _ := range response.Sets {
		elementMap := map[string]interface{}{
			"internal_ip": response.Sets[index].InternalIP,
			"arms_id":     response.Sets[index].ArmsID,
			"app_name":    response.Sets[index].AppName,
			"tree_id":     response.Sets[index].TreeID,
			"idc":         response.Sets[index].Idc,
		}
		DataMap = append(DataMap, elementMap)
	}
	return DataMap
}

//根据类型获取所有的节点Map
func GetNode(nodeType string, userName string) []map[string]interface{} {
	var response ServiceTreeGenerated
	request := gorequest.New().Timeout(10000 * time.Millisecond)

	action := ""
	switch nodeType {
	case "root":
		action = "GetProductNode"
	case "product":
		action = "GetSystemNode"
	case "system":
		action = "GetApplicationNode"
	}

	request.Get(url_config).
		Set("Content-Type", "application/json").
		Send(`{
		"Action":"` + action + `",
		"UserName":"` + userName + `",
		"Project": "Tree",
		"PublicKey": "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com",
		"Signature": "f655e89631c4173ec08438f5b8a70faa841e4a8c"
	}`).
		EndStruct(&response)

	var DataMap = []map[string]interface{}{}
	for index, _ := range response.Data {
		elementMap := map[string]interface{}{
			"app":       response.Data[index].Name,
			"alias":     response.Data[index].Alias,
			"id":        response.Data[index].ID,
			"node_type": response.Data[index].NodeType,
			"parent":    response.Data[index].Parent,
		}
		DataMap = append(DataMap, elementMap)
	}
	return DataMap
}

//获取当前节点的子节点Map
func GetNodeByName(nodeType string, name string) []map[string]interface{} {
	//clear array
	var response ServiceTreeGenerated
	request := gorequest.New().Timeout(10000 * time.Millisecond)

	request.Get(url_config).
		Set("Content-Type", "application/json").
		Send(`{
		"Action":"GetNodeByName",
		"Project": "Tree",
		"PublicKey": "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com",
		"Signature": "f655e89631c4173ec08438f5b8a70faa841e4a8c",
		"NodeType": "` + nodeType + `",
		"Name":"` + name + `"
	}`).
		EndStruct(&response)

	var DataMap = []map[string]interface{}{}
	for index, _ := range response.Data {
		elementMap := map[string]interface{}{
			"app":       response.Data[index].Name,
			"alias":     response.Data[index].Alias,
			"id":        response.Data[index].ID,
			"node_type": response.Data[index].NodeType,
			"parent":    response.Data[index].Parent,
		}
		DataMap = append(DataMap, elementMap)
	}
	return DataMap
}

func timeAlertTips(ctx context.Context, isContextBelowQuery bool) {
	defer wg.Done()

	var progress = 0
	var position = 1
	var seconds = 120
Loop:
	for {
		if progress > 0 {
			fmt.Printf("\033[%dA\033[K", position)
		}
		secondstr := strconv.Itoa(seconds)

		output := fmt.Sprintf(
			"%s%s%s",
			"查询剩余时间倒计时: ",
			SetColor(secondstr, 0, 0, 33),
			" 秒",
		)
		if !isContextBelowQuery {
			fmt.Printf("%s \033[K\n", output)
		}

		seconds--

		if progress >= 120 {
			timeoutMsg := fmt.Sprintf(warningCode, " 当前时间范围内查询超时,请缩小查询时间范围")
			fmt.Println(timeoutMsg)
			break Loop
		}
		progress++
		time.Sleep(time.Duration(1000) * time.Millisecond)

		select {
		case <-ctx.Done():
			break Loop
		default:
		}
	}
}

func QueryRange(filterMap map[string][]string, startTime string, endTime string, step string, limit string, keyword string, direction string, tenant string, isContextBelowQuery bool, user string, env string) QueryRangeGenerate {
	var response QueryRangeGenerate

	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go timeAlertTips(ctx, isContextBelowQuery)

	request := gorequest.New().Timeout(126000 * time.Millisecond)
	filterJsonStr := util.MapToJson(filterMap)

	// if isGrepContext {
	// 	keyword = "nokey"
	// }

	request.Post(url_config).
		// Set("Content-Type", "application/json").
		Param("Action", "QueryRange").
		Param("Project", "Loki").
		Param("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com").
		Param("Signature", "f655e89631c4173ec08438f5b8a70faa841e4a8c").
		Param("Filter", filterJsonStr).
		Param("Direction", direction).
		Param("Start", startTime).
		Param("End", endTime).
		Param("KeyWord", keyword).
		// Param("Step", step).
		Param("Limit", limit).
		Param("Tenant", tenant).
		Param("User", user).
		Param("AppFlag", "shell").
		Param("Env", env).
		EndStruct(&response)

	cancel()
	wg.Wait()

	return response
}

//根据应用名获取日志量等级

func FetchLogScaleByApp(app string, token string) string {
	var response PathGenerated

	request := gorequest.New().Timeout(10000 * time.Millisecond)
	// request.Header.Add("Authorization", "Bearer "+token)
	// fmt.Println("Bearer " + token)
	request.Post(url_config).
		Set("Authorization", "Bearer "+token).
		Set("Content-Type", "application/json").
		Param("Action", "GetLogScaleByApp").
		Param("Project", "Tree").
		Param("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com").
		Param("Signature", "f655e89631c4173ec08438f5b8a70faa841e4a8c").
		Param("App", app).
		EndStruct(&response)

	return response.Data
}

//根据应用获取标签集
func GetLabelSeriesByApp(app string, ip string, starttime string, endtime string, token string, env string) []interface{} {
	// var response LabelsGenerated
	var realData []interface{}

	request := gorequest.New().Timeout(10000 * time.Millisecond)
	_, body, errs := request.Post(url_config).
		Set("Authorization", "Bearer "+token).
		Set("Content-Type", "application/json").
		Param("Action", "GetLabelSeriesByApp").
		Param("Project", "Loki").
		Param("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com").
		Param("Signature", "f655e89631c4173ec08438f5b8a70faa841e4a8c").
		Param("App", app).
		Param("Ip", ip).
		Param("Start", starttime).
		Param("End", endtime).
		Param("Env", env).
		End()
		// EndStruct(&response)
	if len(errs) > 0 {
		promtMsg := fmt.Sprintf(warningCode, "当前标签集数据查询失败")
		fmt.Println(promtMsg)
	} else {
		bodyMap, _ := util.JsonToMap(body)
		resData := bodyMap["Data"].(map[string]interface{})
		realData = resData["data"].([]interface{})
	}
	return realData
}

//根据应用名获取对应租户ID
func GetTenantIdByAppName(app string, token string) string {
	var response PathGenerated

	request := gorequest.New().Timeout(10000 * time.Millisecond)
	// request.Header.Add("Authorization", "Bearer "+token)
	// fmt.Println("Bearer " + token)
	request.Post(url_config).
		Set("Authorization", "Bearer "+token).
		Set("Content-Type", "application/json").
		Param("Action", "GetAppTenant").
		Param("Project", "Loki").
		Param("PublicKey", "9c178e51a7d4dc8aa1dbef0c790b06e7574c4d0ehermestackhui.tu@huifu.com").
		Param("Signature", "f655e89631c4173ec08438f5b8a70faa841e4a8c").
		Param("App", app).
		EndStruct(&response)

	return response.Data
}
