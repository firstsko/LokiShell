package shell

import (
	"fmt"
	"os"
)

func useMsg() {
	usageQueryMore := fmt.Sprintf("过滤查询: more | [grep] [关键字] [2022-01-2812:00:00](开始时间) [2022-01-2812:05:00](结束时间) ")
	usageQueryExampleMore := fmt.Sprintf(warningCode, "Example: more  | grep error 2022-01-2812:00:00 2022-01-2812:05:00 ")
	usageQueryLess := fmt.Sprintf("过滤查询: less | [grep] [关键字] [2022-01-2812:00:00](开始时间) [2022-01-2812:05:00](结束时间) ")
	usageQueryExampleLess := fmt.Sprintf(warningCode, "Example: more  | grep error 2022-01-2812:00:00 2022-01-2812:05:00 ")

	fmt.Fprintln(os.Stdout, usageQueryMore)
	fmt.Fprintln(os.Stdout, usageQueryExampleMore)
	fmt.Fprintln(os.Stdout, usageQueryLess)
	fmt.Fprintln(os.Stdout, usageQueryExampleLess)

	usageTail := fmt.Sprintf("实时查询: tail | [grep] [关键字]")
	usageTailExample := fmt.Sprintf(warningCode, "Example: tail | grep error")
	fmt.Fprintln(os.Stdout, usageTail)
	fmt.Fprintln(os.Stdout, usageTailExample)

}
