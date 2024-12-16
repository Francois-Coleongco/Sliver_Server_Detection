package main

import (
	"app/gatekeeper"
	"app/utils"
	"app/helpers"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)


// URGENT need to make this constantly look at new connecitons. currently it just runs lsof -i once. i need to make it compare against previous. so a diff call on a file might be good i dunno


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

	// CHECK ALSO FOR THE NAMES. check if in your / path NOT THE USER PATH the sudo perm path, for a similar or exact name. sometimes they will attempt to hide themselves via using the name of a legit process such as bash
	// CHECK ALSO FOR THE NAMES. check if in your / path NOT THE USER PATH the sudo perm path, for a similar or exact name. sometimes they will attempt to hide themselves via using the name of a legit process such as bash
	// CHECK ALSO FOR THE NAMES. check if in your / path NOT THE USER PATH the sudo perm path, for a similar or exact name. sometimes they will attempt to hide themselves via using the name of a legit process such as bash

	return executable_path
}

func check_open_files(pid string) []string {
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
	fmt.Println(files_opened_by_pid)
	// gatekeeper.Uses_Shell(files_opened_by_pid)

	return files_opened_by_pid
}

func static_analysis(url_to_executable string) {
	// need to upload the executable to an api of your own so virustotal can receive it

	// OR i needa look more into this, but i think you can hash the executable and you can search for known malwares with that sort of signature
}

func main() {
	// killscore --> uses shell, is obfuscated, uses crypto libs in strings, is tls,

	// 1 means yes | 0 means no

	// kill_score := []byte{0, 0, 0, 0}

	fmt.Println("Current PID:", os.Getpid())

	// find processes making those network connections (lsof)

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("unable to open log file", err)
	}

	defer file.Close()

	// Set the log output to the log file
	log.SetOutput(file)

	conn_lines := lsof_stats()[1:] // start from second line cuz first just gives the column names

	var wg sync.WaitGroup
	pid_chan := make(chan string)

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

			NAME_Field := fields[8] // NAME is an array containing: user, port, address,protocol
			host_name, err := os.Hostname()
			if err != nil {
				log.Println("could not get host_name", err)
			}

			host_name_filter := fmt.Sprintf("%s:", host_name)

			if !strings.Contains(NAME_Field, host_name) {
				fmt.Println("does not contain host")
				continue
			}

			parsed_name_field := strings.Split(NAME_Field, host_name_filter)
			my_port := strings.Split(parsed_name_field[1], "->")

			fmt.Println("PORT", my_port[0])

			fmt.Println(PID_Field)

			executable_path := locate_process(PID_Field)

			fmt.Println("executable_location:", executable_path)

			if _, err := strconv.Atoi(my_port[0]); err == nil {
				wg.Add(2)
				// is a number and can be chucked into the sniffer
				go func() {
					defer wg.Done() // Mark the goroutine as done when it finishes

					utils.Sniffer(my_port[0], PID_Field, pid_chan)
				}()
				go func() {
					defer wg.Done()
					open_files := check_open_files(<-pid_chan)

					uses_shell := gatekeeper.Interacts_With_Shell(open_files)

					if uses_shell {
						// execute strace on that

						children := helpers.Get_Children(<-pid_chan)

						fmt.Println("these are the children", children)

						// children is []string so need to loop through it for tracer in case author tries to spawn a bunch of other seemingly legit child processes
						for i := range children {
							utils.Tracer(children[i]) // this pid is of the imiplant. i need the children. ps --ppid
						}
					}
				}()
			}

		}
	}

	wg.Wait() // some might wait for a ridiculously long time, that's just expected i dont think there's a safer way around this since i need to be able to monitor all network connected processes
}
