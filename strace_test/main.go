package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("strace", "-e", "trace=all", "-p", "99480")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting strace:", err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	errorScanner := bufio.NewScanner(stderr)

	var currentLine string

	go func() {
		for errorScanner.Scan() {
			fmt.Println("strace error:", errorScanner.Text())
		}
	}()

	for scanner.Scan() {
		currentLine = scanner.Text()

		if strings.Contains(currentLine, "write") {
			fmt.Println("Found  syscall:", currentLine)
		}

		fmt.Println("Current line:", currentLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading strace output:", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for strace to finish:", err)
	}
}
