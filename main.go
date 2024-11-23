package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
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

func handle_conn(conn_string string, port_chan chan string) {

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

	port_chan <- local_addr_and_port[len(local_addr_and_port)-1]

}

func find_and_handle_process(port string, pid_chan chan string) {
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

	pid_chan <- strings.Fields(data)[1]

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

func tracer(pid string) {
	fmt.Println("entering tracer")
	// Execute strace with the provided PID
	strace := exec.Command("strace", "-p", pid)

	// Get the stdout pipe of the strace command
	pipe, err := strace.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Start the strace command
	err = strace.Start()
	if err != nil {
		log.Fatalf("Failed to start strace: %v", err)
	}
	defer strace.Wait() // Ensure the strace process finishes when we're done

	// Use a scanner to read the output of the strace command
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		// Print each line of the output
		fmt.Println(scanner.Text())
	}

	// Handle any errors that occurred while scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error while scanning strace output: %v", err)
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

	var wg sync.WaitGroup

	port_chan := make(chan string)

	for i := 0; i < len(conn_lines); i++ {
		if len(conn_lines[i]) > 0 {

			fmt.Println(conn_lines[i])

			wg.Add(1)

			go func() {
				defer wg.Done()
				handle_conn(conn_lines[i], port_chan)
			}()

			// if i run something like tracer anywhere, it will run indefinitely. this is fine, but i need to make sure i can continue searching while the other straces are going. therefore i need go routines here

			// go initial_tracer()

		}
	}

	go func() {
		wg.Wait()
		close(port_chan)
	}()

	pid_chan := make(chan string)

	go func() {
		for port := range port_chan {
			wg.Add(1)
			fmt.Println(port)
			go func() {
				defer wg.Done()
				find_and_handle_process(port, pid_chan)
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(pid_chan)
	}()

	for pid := range pid_chan {
		fmt.Println("huh", pid)
		go tracer(pid)
	}

}
