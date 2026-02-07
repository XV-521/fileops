package util

import (
	"fmt"
	"os"
	"os/exec"
)

func CmdWrap(cmd *exec.Cmd) error {
	cmd.Stdin = nil

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"cmd failed:\npath: %v\nargs: %v\noutput:\n%s",
			cmd.Path,
			cmd.Args,
			string(out),
		)
	}
	return nil
}

// MkdirTempWrap is used to create a temporary directory.
// It accepts a fn and executes it. if the err value returned by fn is nil,
//
// fn: The func that need to be executed after creating the directory.
//
// dir of fn: In MkdirTempWrap, the tmpDir will be passed to fn as dir.
func MkdirTempWrap(fn func(dir string) error) error {
	tmpDir, err := os.MkdirTemp("", "temp-*")
	if err != nil {
		return err
	}

	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	return fn(tmpDir)
}
