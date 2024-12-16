package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

func Get_Children(pid string) []string {
	cmd := exec.Command("ps", "--ppid", pid)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err occurred on getting children of", pid, err)
	}

	// this should be logged womehwere: strings.Split(string(output), "\n") // third and split on space

	// for now just get the pids

	children := strings.Split(string(output), "\n")[1:]

	var child_pids []string

	for i := range children {

		fields := strings.Split(children[i], " ")

		fmt.Println("tnhis is fields", fields)
		child_pids = append(child_pids, fields[1])// starts at 1 because splits on the space between
	}


	return child_pids
}
