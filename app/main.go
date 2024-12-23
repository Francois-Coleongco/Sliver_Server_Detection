package main

import (
	"app/gatekeeper"
	"app/helpers"
	"app/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	//	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// URGENT need to make this constantly look at new connecitons. currently it just runs lsof -i once. i need to make it compare against previous. so a diff call on a file migh be good i dunno

func lsof_stats(lsof_chan chan []string) { // need to make this run on a timer
	cmd := exec.Command("lsof", "-i")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("unable to run `lsof -i`", err)
	}

	str_output := string(output)

	log.Println(str_output)

	conn_lines := strings.Split(str_output, "\n")

	fmt.Println(conn_lines)

	fmt.Println("does this happenaiwghiowajig")
	lsof_chan <- conn_lines
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

func check_open_files(pid string) []string {
	cmd := exec.Command("lsof", "-p", pid)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("could not lsof -p", pid, "| ERROR", err)
	}

	files_opened_by_pid := strings.Split(string(output), "\n")
	fmt.Println(files_opened_by_pid)

	return files_opened_by_pid
}

func static_analysis(url_to_executable string) {
	// need to upload the executable to an api of your own so virustotal can receive it

	// OR i needa look more into this, but i think you can hash the executable and you can search for known malwares with that sort of signature
}

func setup_logs() *log.Logger {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("unable to open log file", err)
	}

	defer file.Close()

	log.SetOutput(file)

	c2_command_log_file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("err at creation/reading of c2_command_log_file")
	}

	c2_command_logger := log.New(c2_command_log_file, "syscall:", log.Ldate)

	return c2_command_logger
}

func main() {
	var wg sync.WaitGroup

	lsof_chan := make(chan []string)

	pid_chan := make(chan string)

	fmt.Println("Current PID:", os.Getpid())

	c2_command_logger := setup_logs()

	fmt.Println("finished setting up logs")

	var conn_lines []string

	go func() {
		conn_lines = <-lsof_chan

		fmt.Println("GOD DAMN IT", conn_lines)
	}()

	lsof_stats(lsof_chan)

	for {
		fmt.Println("execing here")

		fmt.Println(conn_lines)

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

				fmt.Println(PID_Field)

				executable_path := locate_process(PID_Field)

				fmt.Println("executable_location:", executable_path)

				if _, err := strconv.Atoi(my_port[0]); err == nil {
					wg.Add(2)
					// is a number and can be chucked into the sniffer
					go func() {
						defer wg.Done()
						gatekeeper.Strings_Analysis(executable_path)
						utils.Sniffer(my_port[0], PID_Field, pid_chan)
					}()
					go func() {
						defer wg.Done()
						open_files := check_open_files(<-pid_chan)

						uses_shell := gatekeeper.Interacts_With_Shell(open_files)

						if uses_shell {

							// whatever things are being executed in the shell will be logged
							pid := <-pid_chan
							fmt.Println("this is pidchan", pid)
							child_pids := helpers.Get_Children(pid)
							fmt.Println("reached here")
							fmt.Println("these are child_pids", child_pids)

							// children is []string so need to loop through it for tracer in case author tries to spawn a bunch of other seemingly legit child processes
							if len(child_pids) > 0 {
								for i := range child_pids {
									fmt.Println("reached here?")
									if child_pids[i] != "" {
										utils.Tracer(child_pids[i], c2_command_logger)
									}
								}
							}
						}
					}()

				}

			}
		}

		wg.Wait() // some might wait for a ridiculously long time, that's just expected i dont think there's a safer way around this since i need to be able to monitor all network connected processes

		time.Sleep(time.Second)
	}
}
