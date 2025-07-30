package shell

import (
	"fmt"
	"lokishell/kernel"
	"lokishell/util"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

func switchForder(l *readline.Instance, nodeType string, nodeName string, t_path string, direction string) error {

	// update suggestion
	SuggestionMap = kernel.GetNode(nodeType, user)
	SuggestionList = []string{}
	for index, _ := range SuggestionMap {
		SuggestionList = append(SuggestionList, SuggestionMap[index]["app"].(string))
	}

	//switch forder
	if nodeType == "root" {
		CurrentMap = kernel.GetNode("root", user)
	} else if nodeType == "product" {
		CurrentMap = kernel.GetNodeByName(nodeType, nodeName)
	} else if nodeType == "system" {
		CurrentMap = kernel.GetNodeByName(nodeType, nodeName)
	}

	if nodeType == "application" {
		appSlash := strings.Split(t_path, "/")
		app_name = appSlash[3]
		filterMap["_app_"] = []string{app_name}

		instanceMap := kernel.GetInstace(app_name)
		for _, value := range instanceMap {
			ipdb := value["internal_ip"].(string)
			idc := value["idc"].(string)
			ipidcMap[ipdb] = idc
		}

		help(false)

	}

	prefix = "[" + user + "@lokishell " + t_path + "]$"

	l.SetPrompt("\033[32m" + prefix + "\033[0m ")

	s_path = t_path

	return nil
}

func absolute_path(l *readline.Instance, t_path string) bool {
	slash := strings.Split(t_path, "/")
	checkType := len(slash)
	nodeName := slash[len(slash)-1]

	//detect fulldir
	if util.IsContain(fullDirList, t_path) == true {
		if checkType == 2 {
			nodeType = "product"
		} else if checkType == 3 {
			nodeType = "system"
		} else if checkType == 4 {
			nodeType = "application"
		}
		switchForder(l, nodeType, nodeName, t_path, "absolute")
		return false
	} else if t_path == "/" {
		nodeType = "root"
		switchForder(l, nodeType, "/", t_path, "absolute")
		return false
	} else {
		checkmsg := fmt.Sprintf(warningCode, "invalid service tree path,please check!")
		fmt.Fprintln(os.Stdout, checkmsg)
		return false
	}

	return false
}

func cd(l *readline.Instance, t_path string) bool {

	//trim t_path
	t_path = strings.Trim(t_path, " ")

	//detect absolute  or relative
	absolutePath := false
	if strings.HasPrefix(t_path, "/") {
		absolutePath = true
	} else {
		absolutePath = false
	}

	if absolutePath == true {
		//absolute action
		absolute_path(l, t_path)
		return false
	} else {
		//relative action
		//detect suggestionDirList
		//backoff to parent

		if t_path == ".." {
			slash := strings.Split(s_path, "/")
			parentNodeName := slash[:len(slash)-1]
			t_path = strings.Join(parentNodeName, "/")
			if len(slash) == 2 {
				switchForder(l, "root", "/", "/", "absolute")
			} else {
				absolute_path(l, t_path)
			}

			return false
		} else if util.IsContain(SuggestionList, t_path) == true {
			if s_path == "/" {
				t_path = "/" + t_path
			} else {
				t_path = s_path + "/" + t_path
			}
			slash := strings.Split(t_path, "/")
			checkType := len(slash)
			nodeName := slash[len(slash)-1]

			if checkType == 2 {
				nodeType = "product"
			} else if checkType == 3 {
				nodeType = "system"
			} else if checkType == 4 {
				nodeType = "application"
			}
			switchForder(l, nodeType, nodeName, t_path, "relative")
			return false
		} else {
			checkmsg := fmt.Sprintf(warningCode, "Permission error,Please authorize on the service tree,visit http://platform.huifu.com")
			fmt.Fprintln(os.Stdout, checkmsg)
			return false
		}

	}

	return true

}
