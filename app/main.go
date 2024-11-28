package main

import (
	"app/sniff"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func lsof_stats() []string {
	cmd := exec.Command("lsof", "-i")

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("unable to run `lsof -i`", err)
	}

	str_output := string(output)

	log.Println(str_output)

	conn_lines := strings.Split(str_output, "\n")

	return conn_lines

}

func locate_process(pid string) string {

	readlink_args := fmt.Sprintf("/proc/%s/exe", pid)

	cmd := exec.Command("readlink", readlink_args)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("unable to locate process", pid, "| ERROR", err)
	}

	executable_path := string(output)

	return executable_path

}

func sniff_connections(port string) {
	//use the sniffer in the private repo you made to sniff connections from lsof -i

	// plan: use the port to filter packets

	// maybe you can find some unique sliverC2 detections there

	sniffer(port)

}

func check_open_files(pid string) {

	// if user has inotify enabled read from the logs to see network connected processes created or deleted or did something anything with the files on the system

	// how would you do this? listen for changes on the file with diff or something

	// if user does not have inotify enabled then just run file check on the pid

	// lsof -p <pid>

	cmd := exec.Command("lsof", "-p", pid)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("could not lsof -p", pid, "| ERROR", err)
	}

	files_opened_by_pid := strings.Split(string(output), "\n")

}

func static_analysis(url_to_executable string) {

	// need to upload the executable to an api of your own so virustotal can receive it

	// OR i needa look more into this, but i think you can hash the executable and you can search for known malwares with that sort of signature

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

	conn_lines := lsof_stats()[1:] // start from second line cuz first just gives the column names

	for i := 0; i < len(conn_lines); i++ {
		if len(conn_lines[i]) > 0 {

			fmt.Println("new conn_line:", conn_lines[i])

			fields := strings.Fields(conn_lines[i])

			// extract pid from second column

			// COMMAND_Field := fields[0]

			PID_Field := fields[1]

			// USER_Field := fields[2]

			// FD_Field := fields[3]

			// IP_Version := fields[4]

			// DEVICE_Field := fields[5]

			// SIZE_OFF_Field := fields[6]

			// CONN_TYPE_Field := fields[7]

			// NAME := strings.Split(fields[8], ":")

			// NAME is an array containing: user, port, address,protocol

			fmt.Println(PID_Field)

			executable_path := locate_process(PID_Field)

			fmt.Println("executable_location:", executable_path)

		}
	}

}
