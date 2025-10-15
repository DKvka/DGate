package manager

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const (
	LIST = 'l'
	ADD  = 'a'
	DEL  = 'd'
)

var usage string = `
Type the first letter of the command you want to use

[a]dd - Adds a new service to the gateway
[d]elete - Deletes an existing service
[l]ist - Lists every service and port they are running on
`

// Runs a CLI for manual modifications to the gateway
func (m *Manager) RunCli() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		userInput := scanner.Bytes()

		switch bytes.ToLower(userInput)[0] {
		case ADD:
			m.AddService()
		case DEL:
			m.DeleteService()
		default:
			fmt.Println(usage)
		}
	}
}
