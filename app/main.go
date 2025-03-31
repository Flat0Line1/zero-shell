package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func invalidCommand(command string) {
	fmt.Printf("%s: command not found\n", command)
}

func execCommand(args []string) {
	execCommands := make(map[string]string)
	for _, path_to_file := range readPATHDirs() {
		file := filepath.Base(path_to_file)
		if _, exists := execCommands[file]; exists {
			continue
		} else {
			execCommands[file] = path_to_file
		}
	}

	if _, exists := execCommands[args[0]]; exists {
		cmd := exec.Command(args[0], args[1:]...)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		buffer := bufio.NewReader(stdout)
		for {
			line, _, err := buffer.ReadLine()
			if err == io.EOF {
				break
			}
			fmt.Println(string(line))
		}

	} else {
		invalidCommand(args[0])
	}

}

func echoCommand(echo_args []string) {
	fmt.Println(strings.Join(echo_args, " "))
}

func typeCommand(type_arg string) {
	builtin_text := " is a shell builtin"
	func_register := map[string]string{
		"echo": builtin_text,
		"exit": builtin_text,
		"type": builtin_text,
	}

	for _, path_to_file := range readPATHDirs() {
		file := filepath.Base(path_to_file)
		_, exists := func_register[file]
		if exists {
			continue
		} else {
			func_register[file] = " is " + path_to_file
		}
	}

	if _, exists := func_register[type_arg]; exists {
		fmt.Printf("%s%s\n", type_arg, func_register[type_arg])
	} else {
		fmt.Printf("%s: not found\n", type_arg)
	}
}

func readPATHDirs() []string {
	// everything easy change to os.exec.LookPath(command)
	result := []string{}
	pathDirs := os.Getenv("PATH")
	for _, path := range strings.Split(pathDirs, ":") {
		dirContent, _ := os.ReadDir(path)

		for _, entry := range dirContent {
			if entry.IsDir() {
				continue
			}
			result = append(result, path+"/"+entry.Name())
		}
	}

	return result
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
	case "type":
		typeCommand(args[1])
	default:
		execCommand(args)
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
