package gatekeeper

import "strings"

func Uses_Shell(opened_files []string) bool {

	for line := range opened_files {
		if strings.Contains("ptmx", opened_files[line]) {
			return true
		}
	}

	return false

}

func Is_Encrypted(packet string) { // note this function takes in one packet, the for loop should not be in here as it will get really fuckin messy
	// look through packet in the logs to see whether it encrypts the data

	// if the packet shows client hellos and server hellos followed by application data, then it probably is encrypted
}

func Static_Analysis() {

}

func analyze_packet_fmt() { // analyze the packets printed into the sniffy.log

	// check source port

	// check ip

	//check for ACK and ALL OTHER FLAGS OFF, and that the PAYLOAD IS EMPTY [] which will indicate keep alives

	//

}

func striker()
