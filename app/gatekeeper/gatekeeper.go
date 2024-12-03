package gatekeeper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

// the following are helpers for the functions: Is_Encrypted

func strings_analysis(path_to_exec string) {

	// open file and read all printable strings. check if the executable uses crypto libs

	file, err := os.Open(path_to_exec)

	if err != nil {
		log.Println("couldn't open path_to_exec", err)
	}

	reader := bufio.NewReader(file)

	var strings_buf [][]byte

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
			//if it isnt printable then you need to push the current_string onto the strings_buf

			strings_buf = append(strings_buf, current_string)
		}
	}

}

func Interacts_With_Shell(opened_files []string) bool {

	for line := range opened_files {
		if strings.Contains("ptmx", opened_files[line]) {
			return true
		}
	}

	return false

}

func Is_Encrypted(packet string) { // note this function takes in one packet, the for loop should not be in here as it will get really fuckin messy

	strings_analysis()

}

func Static_Analysis() {
	// strings command parsing for encryption libraries used by sliver (i believe they used crypto/tls)
}

func analyze_packet_fmt() { // analyze the packets printed into the sniffy.log

	// check source port

	// check ip

	//check for ACK and ALL OTHER FLAGS OFF, and that the PAYLOAD IS EMPTY [] which will indicate keep alives

	//

}

func striker() // function to tally a score on whether process should be killed
