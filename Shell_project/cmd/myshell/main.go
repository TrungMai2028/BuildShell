package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	for {
		// Prompt user to type in
		fmt.Fprint(os.Stdout, "$ ")

		// Get user input from keyboard
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		// Handle errors
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Split the input to separate the command and the arguments
		cmd := strings.Split(strings.TrimSpace(input), " ")

		// Command and arguments
		command := cmd[0]
		args := cmd[1:]

		// Handle commands
		switch command {
		case "exit":
			handleExit(args)
		case "echo":
			handleEcho(args)
		case "type":
			handleType(args)

		case "pwd":
			handlePwd()

		case "cd":
			handleCd(args)

		case "clear":
			handleClear()

		case "dir": //dir for Window
			command := exec.Command("cmd", "/c", "dir /b")

			// cmd: Use the default shell (cmd.exe on Windows)
			// /c: Execute the following command and exit
			// dir /b: List directory contents in bare format

			command.Stderr = os.Stderr
			command.Stdout = os.Stdout
			err := command.Run()
			if err != nil {
				fmt.Println("Error running dir:", err)

			}

		default: //for linux and Mac
			command1 := exec.Command(command, args...)
			command1.Stderr = os.Stderr
			command1.Stdout = os.Stdout

			if err := command1.Run(); err != nil {
				fmt.Printf("%s: command not found\n", command)
			}

		}
	}
}

// handleExit processes the "exit" command
func handleExit(args []string) {
	// Default exit code is 0
	code := 0
	if len(args) > 0 {
		// Parse the exit code if provided
		parsedCode, err := strconv.Atoi(args[0])
		if err != nil {
			code = 1
		} else {
			code = parsedCode
		}
	}
	// Exit the program with the specified code
	os.Exit(code)
}

// handleEcho processes the "echo" command
func handleEcho(args []string) {
	// Print the arguments joined by a space
	fmt.Println(strings.Join(args, " "))
}

// handleType processes the "type" command
func handleType(args []string) {
	if len(args) < 1 {
		// If no argument is provided, print an error message
		fmt.Println("type: no argument provided")
		return
	}

	arg := args[0]
	// Check if the argument is a shell builtin
	if isShellBuiltin(arg) {
		fmt.Printf("%s is a shell builtin\n", arg)
	} else {
		// Check for the executable in the PATH
		checkExecutable(arg)
	}
}

// isShellBuiltin checks if a command is a shell builtin
func isShellBuiltin(arg string) bool {
	switch arg {
	case "exit", "echo", "type", "pwd":
		return true
	default:
		return false
	}
}

// checkExecutable searches for the executable in the PATH
func checkExecutable(arg string) {
	// Get the PATH environment variable
	env := os.Getenv("PATH")
	// Split PATH into directories
	paths := strings.Split(env, string(os.PathListSeparator))
	found := false
	// Check for common Windows executable extensions
	extensions := []string{"", ".exe", ".bat", ".cmd"}

	// Iterate through each directory in PATH
	for _, path := range paths {
		// Check for the executable with each extension
		for _, ext := range extensions {
			fp := filepath.Join(path, arg+ext)
			// If the file exists, print the path and set found to true
			if _, err := os.Stat(fp); err == nil {
				fmt.Printf("%v is %v\n", arg, fp)
				found = true
				break
			}
		}
		// If the executable is found, break out of the loop
		if found {
			break
		}
	}
	// If the executable is not found, print an error message
	if !found {
		fmt.Printf("%s: not found\n", arg)
	}
}

// handleType processes the "pwd" command
func handlePwd() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(path)
}

// handleType processes the "cd" command
func handleCd(args []string) {
	if len(args) == 0 {
		fmt.Println("Please specify a directory")
	}

	dir := args[0]
	if dir == "~" {
		homeDir := os.Getenv("USERPROFILE")
		if homeDir != "" {
			dir = homeDir
		} else {
			return
		}
	}
	//change new directory
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
		return
	}

}

// handleType processes the "clear" command
func handleClear() {
	var cmd *exec.Cmd
	//runtime.GOOS represents the operating system Go program is running
	switch runtime.GOOS {
	case "windows": // For Windows, create a command to run "cmd /c cls"
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
