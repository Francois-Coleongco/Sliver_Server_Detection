package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func lsof_stats() []string {
	cmd := exec.Command("lsof -i")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("unable to run ss", err)
	}

	conn_lines := strings.Split(string(output), "\n")

	return conn_lines

}

func sniff_connections(port string) {
	//use the sniffer in the private repo you made to sniff connections from lsof -i

	// maybe you can find some unique sliverC2 detections there
}

func check_open_files() {

	// if user has inotify enabled read from the logs to see network connected processes created or deleted or did something anything with the files on the system

	// how would you do this? listen for changes on the file with diff or something

	// if user does not have inotify enabled then just run file check on the pid

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
			}
		}
	}

}
