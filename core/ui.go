package core

import (
	"fmt"
	"os/exec"
)

// java -jar simone.jar
func openSimonUI() error {
	path, err := exec.LookPath("java")
	if err != nil {
		return fmt.Errorf("lookup java failed:%v", err)
	}
	cmd := exec.Command(path, "-jar", "simone.jar")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("run java failed (missing simone.jar?) :%v", err)
	}
	return nil
}
