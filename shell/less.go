package shell

import (
	"github.com/chzyer/readline"
)

func less(l *readline.Instance, filterMap map[string][]string, starttime string, endtime string, step string, limit string, keyword string, grep bool, timepoint bool) bool {

	fetchEvent(filterMap, starttime, endtime, step, limit, keyword, "BACKWARD", true, false)

	return true

}
