package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")
		user_command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		user_command = strings.TrimSpace(user_command)
		switch user_command {
		case "exit 0":
			os.Exit(0)
		default:
			fmt.Printf("%s: command not found\n", user_command)
		}
	}

}
