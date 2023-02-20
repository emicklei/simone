package api

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type ProgramContext struct {
}

// OnLogin immediately askes for the credentials to setup the plugin
func (p ProgramContext) OnLogin(plugin Plugin, logincall LoginFunc) {
	usr, pwd, err := credentials(plugin.Namespace())
	if err != nil {
		log.Println("failed to get credentials for ", plugin.Namespace())
		return
	}
	if err := logincall(usr, pwd); err != nil {
		log.Println("failed to login for ", plugin.Namespace())
	}
}

// Set is a no-op for programs
func (p ProgramContext) Set(name string, value any) error { return nil }

func credentials(who string) (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("[%s] username: ", who)
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Printf("[%s] password: ", who)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
