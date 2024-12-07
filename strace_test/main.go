package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Define the command to run and strace it
	cmd := exec.Command("strace", "-p", "142936")

	// Set up the command to capture stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished running strace")
}
