package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func invalidCommand(command string) {
	fmt.Printf("%s: command not found\n", command)
}

func echoCommand(echo_args []string) {
	fmt.Println(strings.Join(echo_args, " "))
}

func handler(command string) {
	command = strings.TrimSpace(command)
	var args = strings.Split(command, " ")

	switch args[0] {
	case "exit":
		if len(args) > 1 {
			exit_code, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid exit code")
				return
			}
			os.Exit(exit_code)
		}
		return
	case "echo":
		echoCommand(args[1:])
	default:
		invalidCommand(command)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")
		user_command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		handler(user_command)
	}

}
