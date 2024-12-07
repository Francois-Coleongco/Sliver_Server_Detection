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

func packet_is_encrypted(packet string) { // note this function takes in one packet, the for loop should not be in here as it will get really fuckin messy

	// okay so in order to make sure that the data is encrypted using TLS, we need to know if there was a key exchange that went on. to do that we would need to sniff the packets just when the implant begins communication.

	// so we need to look at some application data packet on our logs, and see if there was a key exchange that occurred before it to see if there was an actual TLS encryption that occurred.

	// step 1: check if packet contains

}

func Static_Analysis(path_to_exec string) {
	// strings command parsing for encryption libraries used by sliver (i believe they used crypto/tls)

	strings_analysis(path_to_exec)
}

func analyze_packet_fmt() {
	// check source port

	// check ip

	//check for ACK and ALL OTHER FLAGS OFF, and that the PAYLOAD IS EMPTY [] which will indicate keep alives

	//

}
