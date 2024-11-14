package main

import (
	"bufio"
	"fmt"
	"io"
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

			fmt.Println("handle_conn ====>", port)
			fmt.Println("find_and_handle_process ====>", pid)
			pipe := tracer(pid) // this just constructs the command and returns the pipe, no need for go routine here

			// if i run something like tracer anywhere, it will run indefinitely. this is fine, but i need to make sure i can continue searching while the other straces are going. therefore i need go routines here

			// need goroutine for initial_tracer

			go initial_tracer()

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

func recon(pid string) {

	ps_cmd := exec.Command("watch -n 1", "ps", "--ppid", pid)

	pipe, err := ps_cmd.StdoutPipe()

	if err != nil {
		log.Fatal("unable to retrieve stdout pipe from ps --ppid", err)
	}

	err = ps_cmd.Start()

	if err != nil {
		log.Fatal("ps --ppid could not start", err)
	}

	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println("from ps --ppid ====>", line)

	}

}
func tracer(pid string) io.ReadCloser {
	strace := exec.Command("strace", "-p", pid)

	pipe, err := strace.StdoutPipe()

	if err != nil {
		log.Fatal("unable to retrieve stdout pipe from strace", err)
	}

	err = strace.Start()

	if err != nil {
		log.Fatal("strace could not start", err)
	}

	return pipe

}

func initial_tracer(pipe io.ReadCloser) bool {
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println("from strace:", line)

		enc_comms := strings.Contains(line, "\\27\\3\\3\\")

		if enc_comms {
			// when you see this, you should immediately run ps --ppid <current_pid> to see any child processes spawned by the one you're currently looking at
			return true
		} // omg i almost made this return false./ that woulda just completely killed the for loop :sob:
	}

	return false

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
