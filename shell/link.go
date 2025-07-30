package shell

import (
	"fmt"
	"lokishell/kernel"
	"os"

	"github.com/chzyer/readline"
)

func link(l *readline.Instance, app string) bool {

	app_name = app
	filterMap["_app_"] = []string{app_name}
	instanceMap := kernel.GetInstace(app_name)
	for _, value := range instanceMap {
		ipdb := value["internal_ip"].(string)
		idc := value["idc"].(string)
		ipidcMap[ipdb] = idc
	}

	s_path = kernel.FetchPathByApp(app_name)
	nodeType = "application"

	moreCompleter.SetChildren(
		makeMoreCompletion(app_name),
	)
	lsCompleter.SetChildren(
		makeLsCompletion(app_name),
	)
	l.Config.AutoComplete = completer

	if s_path == "/invalid app" {
		checkmsg := fmt.Sprintf(warningCode, "invalid service app,please check!")
		fmt.Fprintln(os.Stdout, checkmsg)
		return false
	} else {
		help(false)
		prefix = "[" + user + "@lokishell " + s_path + "]$"
		l.SetPrompt("\033[32m" + prefix + "\033[0m ")

	}

	return true
}
