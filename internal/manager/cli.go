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
			name := GetString("Enter the name of the service")
			addr := GetString("Address of the service")
			m.AddService(name, addr)

		case DEL:
			name := GetString("Enter the name of the service to delete or leave blank to cancel")
			m.DeleteService(name)

		default:
			fmt.Println(usage)
		}
	}
}

func GetString(prompt string) string {
	fmt.Println(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}
