package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func get_strings(path_to_exec string) ([]byte, []string) {

	// open file and read all printable strings. check if the executable uses crypto libs

	var status []byte

	file, err := os.Open(path_to_exec)

	if err != nil {
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var strings_buf []string
	var current_string []byte

	for {
		char, err := reader.ReadByte()

		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading file:", err)
			}
			break
		}

		if unicode.IsPrint(rune(char)) {
			current_string = append(current_string, char)
		} else {
			if len(current_string) > 0 {
				current_string = append(current_string, '\n')
				strings_buf = append(strings_buf, string(current_string))
				current_string = nil
			}
		}
	}

	if len(current_string) > 0 {
		strings_buf = append(strings_buf, string(current_string))
	}

    if uses crypto stuff then {
        status = append(status, 1)
    }

	return strings_buf

}


func obfuscation_check(strings_buf string) {

    // SHOULD BE INVOKED BEFORE strings_analysis BUT AFTER get_strings

    // check the strings for gibberish looking obfuscation. you could use the oh so shiny and cool nlp model wooooooo

}

func strings_analysis(strings_buf []string) {
    // look for things related to encryption in the libraries output
}

func main() {
	fmt.Println(get_strings("/home/hitori/kodoku/Sliver_Server_Detection/research/testing_sniff/main"))
}
