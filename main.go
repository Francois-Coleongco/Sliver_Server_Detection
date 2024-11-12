package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func handle_conn(conn_string string) string {

	conn_string_fields := strings.Fields(conn_string)

	conn_type := conn_string_fields[0]

	conn_status := conn_string_fields[1]

	receive_queue := conn_string_fields[2]

	send_queue := conn_string_fields[3]

	local_addr := conn_string_fields[4]

	remote_addr := conn_string_fields[5]

	log.Println(conn_type, conn_status, receive_queue, send_queue, local_addr, remote_addr)

	local_addr_and_port := strings.Split(local_addr, ":")

	return local_addr_and_port[1]

}

func find_and_handle_process(port string) {
	// execute lsof here

	// grab ALL THE INFORMATION IN THERE

}

func main() {

	// find network connections (ss)
	// find processes making those network connections (lsof)
	// use the pid found in previous step to strace the process
	// in strace output, see if it makes encrypted communications (tls)
	//

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err) // If the file can't be opened, exit the program
	}
	defer file.Close()

	// Set the log output to the log file
	log.SetOutput(file)
	ss_run := exec.Command("ss")

	ss_output, err := ss_run.CombinedOutput()

	ss_run_filter_tcp := exec.Command("grep", "tcp")

	ss_run_filter_tcp.Stdin = strings.NewReader(string(ss_output))

	if err != nil {
		log.Fatalf("ss did not work!!!")
	}
	output, err := ss_run_filter_tcp.CombinedOutput()

	conn_lines := strings.Split(string(output), "\n")

	for i := 0; i < len(conn_lines); i++ {
		handle_conn(conn_lines[i])
	}

}
