package shell

import (
	"bufio"
	"fmt"
	"lokishell/kernel"
	"os"
	"strings"
	"syscall"

	"github.com/chzyer/readline"
	"golang.org/x/term"
)

func login(l *readline.Instance, userAlias string) {
	if userAlias == "nouser" {
		fmt.Println("please input Username:")
		reader := bufio.NewReader(os.Stdin)
		username, _ := reader.ReadString('\n')
		fmt.Println("please input Password:")
		bytepw, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		pass := string(bytepw)
		user = strings.TrimSpace(username)
		fmt.Println("Login with " + user + "password:" + pass)

		token = kernel.FetchTokenByUsername(user)

		prefix = "[" + user + "@lokishell> ~]$"
	} else {
		prefix = "[" + user + "@lokishell> ~]$"
	}

	l.SetPrompt("\033[32m" + prefix + "\033[0m ")

}
