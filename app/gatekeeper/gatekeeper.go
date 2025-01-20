package gatekeeper

import (
	"app/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func Strings_Analysis(path_to_exec string) {
	// open file and read all printable strings. check if the executable uses crypto libs
	file, err := os.Open(path_to_exec)
	if err != nil {
		log.Println("couldn't open path_to_exec", err)
	}

	reader := bufio.NewReader(file)

	var current_string []byte

	for {
		char, err := reader.ReadByte()
		if err != nil {
			// If EOF reached, break the loop
			if err.Error() != "EOF" {
				fmt.Println("end of file!")
			}
			break
		}

		if unicode.IsPrint(rune(char)) {
			// Accumulate printable characters
			current_string = append(current_string, char)
		} else {
			// if it isnt printable then you need to push the current_string onto the strings_buf

			// note this is where a string is complete

			go func() {
				utils.Check_For_Crypto_Libs(string(current_string))
			}()
		}
	}
}

func Interacts_With_Shell(opened_files []string) bool { // if it interacts with the shell, attempt to see the strace of the /dev/pts/*

	for line := range opened_files {
		if strings.Contains("ptmx", opened_files[line]) {
			return true
		}
	}

	return false
}


func Process_Killer(pid string) bool {

	cmd := exec.Command("kill", pid)

	exit_code := cmd.ProcessState.ExitCode()
	
	return exit_code == 0 // if it exited with a successful kill, return true. else return false
}
