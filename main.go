package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func socket_stats() {
	ss_run := exec.Command("ss")

	ss_output, err := ss_run.CombinedOutput()

	if err != nil {
		log.Fatal("unable to run ss", err)
	}

	ss_run_filter_tcp := exec.Command("grep", "tcp")

	ss_run_filter_tcp.Stdin = strings.NewReader(string(ss_output))

	output, err := ss_run_filter_tcp.CombinedOutput()

	conn_lines := strings.Split(string(output), "\n")

	for i := 0; i < len(conn_lines); i++ {
		if len(conn_lines[i]) > 0 {

			fmt.Println(conn_lines[i])

			port := handle_conn(conn_lines[i])
			pid := find_and_handle_process(port)

			fmt.Println("this is from handle_conn", port)
			fmt.Println("this is from find_and_handle_process")
			recon(pid)
		}

	}
}

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

	return local_addr_and_port[len(local_addr_and_port)-1]

}

func find_and_handle_process(port string) string {
	// execute lsof here
	lsof_cmd_args := fmt.Sprintf(":%s", port)

	cmd := exec.Command("lsof", "-i", lsof_cmd_args)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("could not run lsof", err)
		fmt.Println(err)
	}

	string_out := string(output)

	output_arr := strings.Split(string_out, "\n")

	// titles := output_arr[0]
	/* COMMAND   PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
	 */

	data := output_arr[1]

	/* brave   36615  soy  106u  IPv4  61877      0t0  TCP pop-os:35774->151.101.21.140:https (ESTABLISHED) */

	fmt.Println("this is pid", strings.Fields(data)[1])

	return strings.Fields(data)[1]

}

func recon(pid string) string {

	ps_cmd := exec.Command("ps", "aux")

	output, err := ps_cmd.CombinedOutput()

	if err != nil {
		log.Fatal("unable to run ps", err)
	}

	ps_filter := exec.Command("grep", pid)

	ps_filter.Stdin = strings.NewReader(string(output))

	ps_filtered_output, err := ps_filter.CombinedOutput()

	ps_filtered_output_str := string(ps_filtered_output)
	fmt.Println("from ps", ps_filtered_output_str)

	log.Println(ps_filtered_output_str)
	return ps_filtered_output_str

}

func tracer(pid string) {
	strace := exec.Command("strace", "-p", pid)

	pipe, err := strace.StdoutPipe()

	if err != nil {
		log.Fatal("unable to retrieve stdout pipe from strace", err)
	}

	err = strace.Start()

	if err != nil {
		log.Fatal("strace could not start", err)
	}

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "i forgot what it looked like but it had write and \\27\\3\\3\\etc etc i dont remember. i think that was it though") {
			fmt.Println("EVILLLL")
		}
	}

}

func main() {

	// find network connections (ss)
	// find processes making those network connections (lsof)
	// use the pid found in previous step to strace the process
	// in strace output, see if it makes encrypted communications (tls)

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("unable to open log file", err)
	}
	defer file.Close()

	// Set the log output to the log file
	log.SetOutput(file)

	socket_stats()

}
