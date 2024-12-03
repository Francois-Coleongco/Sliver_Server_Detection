package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func strings_analysis(path_to_exec string) {

	// open file and read all printable strings. check if the executable uses crypto libs

	file, err := os.Open(path_to_exec)

	if err != nil {
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var strings_buf [][]byte

	var current_string []byte

	for {
		fmt.Println("checking")
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

	fmt.Println(strings_buf)

}

func main() {
	strings_analysis("/home/hitori/kodoku/Sliver_Server_Detection/research/testing_sniff/main")
}
