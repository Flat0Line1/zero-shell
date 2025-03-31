package main

import (
	"bufio"
	"errors"
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

func pwdCommand() {
	// currDir, _ := os.Getwd()
	// fmt.Println(currDir)
	fmt.Println(os.Getenv("PWD"))
}

func cdCommand(args []string) {
	path := strings.TrimSpace(args[0])
	isAbsPath := path[0] == '/'

	if path == "~" {
		path = os.Getenv("HOME")
	} else if !isAbsPath {
		path = filepath.Join(os.Getenv("PWD"), path)
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s: No such file or directory\n", path)
		return
	}
	os.Setenv("PWD", path)
	// if err := os.Chdir(command); err != nil {
	// 	fmt.Fprintf(os.Stdout, "%s: No such file or directory\n", command)
	// }
}

func typeCommand(type_arg string) {
	builtin_text := " is a shell builtin"
	func_register := map[string]string{
		"echo": builtin_text,
		"exit": builtin_text,
		"type": builtin_text,
		"pwd":  builtin_text,
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

func parseArgs(command string) []string {
	var args []string
	var current_tokens strings.Builder
	inSingleQuotes, inDoubleQuotes, isEscapedSeq := false, false, false

	for _, char := range command {
		if isEscapedSeq {
			current_tokens.WriteRune(char)
			isEscapedSeq = false
			continue
		}

		switch char {
		case '\\':
			isEscapedSeq = true
		case '"':
			inDoubleQuotes = !inDoubleQuotes && !inSingleQuotes
			if inSingleQuotes {
				current_tokens.WriteRune('"')
			}
		case '\'':
			inSingleQuotes = !inSingleQuotes && !inDoubleQuotes
			if inDoubleQuotes {
				current_tokens.WriteRune('\'')
			}
		case ' ':
			if inSingleQuotes || inDoubleQuotes {
				current_tokens.WriteRune(' ')
				continue
			}
			if current_tokens.Len() != 0 {
				args = append(args, current_tokens.String())
				current_tokens.Reset()
			}
		default:
			current_tokens.WriteRune(char)
		}

	}
	if current_tokens.Len() > 0 {
		args = append(args, current_tokens.String())

	}

	return args
}

func handler(command string) {
	command = strings.Trim(command, "\r\n")
	args := parseArgs(command)

	switch args[0] {
	case "exit":
		exit_code := 0
		if len(args) > 1 {
			new_code, err := strconv.Atoi(args[1])
			exit_code = new_code
			if err != nil {
				fmt.Println("Invalid exit code")
				return
			}
		}
		os.Exit(exit_code)
	case "echo":
		echoCommand(args[1:])
	case "pwd":
		pwdCommand()
	case "cd":
		cdCommand(args[1:])
	case "type":
		typeCommand(args[1])
	default:
		execCommand(args)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		// currDir, _ := os.Getwd()
		// splitedCurrDir := strings.Split(currDir, "/")
		// lastThreeEntries := strings.Join(splitedCurrDir[len(splitedCurrDir)-1:], "/")
		// fmt.Fprint(os.Stdout, lastThreeEntries+" $ ")
		fmt.Fprint(os.Stdout, "$ ")
		user_command, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		handler(user_command)
	}

}
