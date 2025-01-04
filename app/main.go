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

	lsof_chan <- conn_lines
}

func locate_process(pid string, pids_in_processing *map[string]struct{}) string {
	readlink_args := fmt.Sprintf("/proc/%s/exe", pid)

	cmd := exec.Command("readlink", readlink_args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("unable to locate process", pid, "| ERROR", err)
		// remove from pids_in_processing
		delete(*pids_in_processing, pid)
		fmt.Println("deleting pid in process")
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

	return files_opened_by_pid
}

func static_analysis(url_to_executable string) {
	// need to upload the executable to an api of your own so virustotal can receive it

	// OR i needa look more into this, but i think you can hash the executable and you can search for known malwares with that sort of signature
}

func setup_c2_logs(c2_file_id string) *log.Logger {
	err := os.MkdirAll("./c2_logs", 0666)
	if err != nil {
		fmt.Println("couldn't create c2_logs directory")
	}

	file_name := fmt.Sprintf("./c2_logs/c2_log_%s.log", c2_file_id)

	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("unable to open log file", err)
	}

	defer file.Close()

	c2_command_logger := log.New(file, "syscall:", log.Ldate)

	return c2_command_logger
}

func entry(wg *sync.WaitGroup, pid_chan chan string, lsof_chan chan []string, pids_in_processing *map[string]struct{}, counter *int) {
	fmt.Printf("CURRENTLY RUNNING ENTRY WITH COUNTER %v\n", (*counter))
	lines := <-lsof_chan


	for i := 0; i < len(lines); i++ {
		if len(lines[i]) > 0 {


			fields := strings.Fields(lines[i])

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

			executable_path := locate_process(PID_Field, pids_in_processing)

			fmt.Println("executable_location:", executable_path)

			if _, err := strconv.Atoi(my_port[0]); err == nil {
				wg.Add(2)
				// is a number and can be chucked into the sniffer
				go func() {
					defer wg.Done()

					pid := <-pid_chan

					if _, exists := (*pids_in_processing)[pid]; exists {
						fmt.Printf("this pid %s is still being processed\n", pid)
						return // exit go func since we dont wanna double process on this pid
					}

					(*pids_in_processing)[pid] = struct{}{}

					open_files := check_open_files(pid)

					uses_shell := gatekeeper.Interacts_With_Shell(open_files)

					if uses_shell {

						// whatever things are being executed in the shell will be logged
						fmt.Println("this is pidchan", pid)
						child_pids := helpers.Get_Children(pid)
						fmt.Println("reached here")
						fmt.Println("these are child_pids", child_pids)

						// children is []string so need to loop through it for tracer in case author tries to spawn a bunch of other seemingly legit child processes
						if len(child_pids) > 0 {
							for i := range child_pids {
								fmt.Println("reached here?")
								if child_pids[i] != "" {
									c2_command_logger := setup_c2_logs(pid)
									utils.Tracer(child_pids[i], c2_command_logger)
								}
							}
						}
					}
				}()
				go func() {
					defer wg.Done()
					gatekeeper.Strings_Analysis(executable_path)
					utils.Sniffer(my_port[0], PID_Field, pid_chan)
					fmt.Println("does it get past snniffer")
				}()

			}

		}
	}

	wg.Wait() // some might wait for a ridiculously long time, that's just expected i dont think there's a safer way around this since i need to be able to monitor all network connected processes

	time.Sleep(time.Second)
}

func main() {
	main_log_file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("couldn't open log file app.log")
	}

	log.SetOutput(main_log_file)

	var wg sync.WaitGroup

	lsof_chan := make(chan []string)

	pid_chan := make(chan string)

	pids_in_processing := make(map[string]struct{})

	fmt.Println("Current PID:", os.Getpid())

	fmt.Println("finished setting up logs")

	counter := 0

	var counter_pointer *int = &counter

	for {
		go func() {
			fmt.Println("new lsof routine")
			fmt.Println("new lsof routine")
			fmt.Println("new lsof routine")
			fmt.Println("new lsof routine")
			lsof_stats(lsof_chan)
			time.Sleep(time.Second * 5)
		}()

		go func() {
			fmt.Println("new entry routine")
			fmt.Println("new entry routine")
			fmt.Println("new entry routine")
			fmt.Println("new entry routine")
			entry(&wg, pid_chan, lsof_chan, &pids_in_processing, &counter)
			*counter_pointer += 1
		}()

		time.Sleep(time.Second * 5)

	}
}
