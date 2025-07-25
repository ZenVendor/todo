package main

import (
	"os"
	"os/exec"
)

func (c *Config) GetTextFromEditor(current string) (desc string, err error) {
	tempFile, err := os.CreateTemp("", "input*.txt")
	if err != nil {
		return current, err
	}
	defer tempFile.Close()

	if current != "" {
		tempFile.WriteString(current)
		tempFile.Sync()
	}

	editor := c.Editor
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return current, err
	}
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return current, err
	}
	return string(content), err
}
