package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	for {
		fmt.Fprint(os.Stdout, "$ ")
		user_command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(user_command[:len(user_command)-1] + ": command not found")
	}
}
