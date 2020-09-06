package common

import (
	"os/exec"
)

var Command command

type command struct {
}

func (c *command) CopyFile(sourceFilePath string, targetFilePath string) error {
	return exec.Command("cmd", "/C", "copy", sourceFilePath, targetFilePath, "/y").Run()
}

func (c *command) MakeDirectory(directoryPath string) error {
	if err := exec.Command("cmd", "/C", "mkdir", directoryPath).Run(); nil != err {
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}

		return err
	}

	return nil
}
