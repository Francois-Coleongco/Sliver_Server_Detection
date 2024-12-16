package helpers

import (
	"os/exec"
	"fmt"
	"strings"
)

func Get_Children(pid string) []string {

	cmd := exec.Command("ps", "--ppid", pid)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("err occurred on getting children of", pid, err)
	}

	return strings.Split(string(output), "\n") // third and split on space
}
