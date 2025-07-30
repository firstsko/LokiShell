package shell

import (
	"fmt"
)

func banner() {
	//Show Banner
	msg1 := fmt.Sprintf(successCode, ` _     ___  _  _____ ____  _   _ _____ _     _`)
	msg2 := fmt.Sprintf(successCode, `| |   / _ \| |/ /_ _/ ___|| | | | ____| |   | |`)
	msg3 := fmt.Sprintf(successCode, `| |  | | | | ' / | |\___ \| |_| |  _| | |   | |`)
	msg4 := fmt.Sprintf(successCode, `| |__| |_| | . \ | | ___) |  _  | |___| |___| |___`)
	msg5 := fmt.Sprintf(successCode, `|_____\___/|_|\_\___|____/|_| |_|_____|_____|_____|`)
	msg6 := fmt.Sprintf(infoCode, `┌  Next Generation Log Center ───────────────────────────────────────────────┐`)
	msg7 := fmt.Sprintf(infoCode, `│                                                                            │`)
	msg8 := fmt.Sprintf(infoCode, `│   Schema   Based on ARMS Product-System-Application CMDB                   │`)
	msg9 := fmt.Sprintf(infoCode, `│                                                                            │`)
	msg10 := fmt.Sprintf(infoCode, `│   Filter   Compatible with Linux Bash commands                             |`)
	msg11 := fmt.Sprintf(infoCode, `│                                                                            │`)
	msg12 := fmt.Sprintf(infoCode, `│   Search  Query Keyword by Syntax                                          │`)
	msg13 := fmt.Sprintf(infoCode, `└────────────────────────────────────────────────────────────────────────────┘`)

	fmt.Println(msg1)
	fmt.Println(msg2)
	fmt.Println(msg3)
	fmt.Println(msg4)
	fmt.Println(msg5)
	fmt.Print("\n")
	fmt.Println(msg6)
	fmt.Println(msg7)
	fmt.Println(msg8)
	fmt.Println(msg9)
	fmt.Println(msg10)
	fmt.Println(msg11)
	fmt.Println(msg12)
	fmt.Println(msg13)
	fmt.Println("输入login登陆账户,输入help查看帮助手册")
}
