package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func socket_stats() []string {
	ss_run := exec.Command("ss")

	ss_output, err := ss_run.CombinedOutput()

	if err != nil {
		log.Fatal("unable to run ss", err)
	}

	ss_run_filter_tcp := exec.Command("grep", "tcp")

	ss_run_filter_tcp.Stdin = strings.NewReader(string(ss_output))

	output, err := ss_run_filter_tcp.CombinedOutput()

	conn_lines := strings.Split(string(output), "\n")

	return conn_lines

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

	fmt.Println(local_addr_and_port[len(local_addr_and_port)-1], "DEBUG")

	return local_addr_and_port[len(local_addr_and_port)-1]

}

func find_and_handle_process(port string) string {
	// execute lsof here lsof_cmd_args := fmt.Sprintf(":%s", port)

	lsof_cmd_args := fmt.Sprintf(":%s", port)

	fmt.Println("this is port", port)

	cmd := exec.Command("lsof", "-i", lsof_cmd_args)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("could not run lsof", err)
		fmt.Println("here", err)
	}

	string_out := string(output)
	print("here", string_out)

	if len(string_out) > 0 {
		output_arr := strings.Split(string_out, "\n")

		// titles := output_arr[0]
		/* COMMAND   PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
		 */

		data := output_arr[1]

		/* brave   36615  soy  106u  IPv4  61877      0t0  TCP pop-os:35774->151.101.21.140:https (ESTABLISHED) */
		fmt.Println("this is pid", strings.Fields(data)[1])
		return strings.Fields(data)[1]
	}

	return ""
}

func ps_recon(pid string) string {

	ps_cmd := exec.Command("ps", "--ppid", pid) // listen for children | the sliver implant may spawn another process for a shell

	output, err := ps_cmd.CombinedOutput()

	if err != nil {
		log.Println("couldnt ps", err)
	}

	return string(output)

}

func tracer(pid string) {
	cmd := exec.Command("/bin/sh", "-c", "sudo strace -p 3515")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting StdoutPipe:", err)
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error getting StderrPipe:", err)
		return
	}

	// Create a reader to capture both stdout and stderr
	go func() {
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stdout:", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text()) // Print to stderr
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stderr:", err)
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command to finish:", err)
	}

}

func main() {
	fmt.Println("Current PID:", os.Getpid())

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

	conn_lines := socket_stats()

	for i := 0; i < len(conn_lines); i++ {
		if len(conn_lines[i]) > 0 {

			fmt.Println(conn_lines[i])

			port := handle_conn(conn_lines[i])

			pid := find_and_handle_process(port)

			if pid != "" { // handling no pid case
				fmt.Println(ps_recon(pid))
				tracer(pid)
			}
		}
	}

}
