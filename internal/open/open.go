package open

import (
	"os/exec"
)

func InEditor(editor string, path string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

func InMacOS(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Start()
}
