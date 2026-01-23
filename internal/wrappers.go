package internal

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func CmdWrapper(cmd *exec.Cmd) error {
	cmd.Stdin = nil

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed: %v", string(out))
		return err
	}
	return nil
}

// MkdirTempWrapper is used to create a temporary directory.
// It accepts a fn and executes it. if the err value returned by fn is nil,
// the temporary directory will be moved to finalPath.
//
// finalPath: The path to the final dir.
//
// fn: The func that need to be executed after creating the directory.
//
// dir of fn: In MkdirTempWrapper, the tmpPath will be passed to fn as dir.
func MkdirTempWrapper(finalPath string, fn func(dir string) error) error {
	tmpPath, err := os.MkdirTemp("", "temp-*")
	if err != nil {
		return err
	}

	defer func() {
		_ = os.RemoveAll(tmpPath)
	}()

	err = fn(tmpPath)
	if err != nil {
		return err
	}

	return os.Rename(tmpPath, finalPath)
}

type BatchMode struct {
	Async  bool
	Strict bool
}

func DoBatchWrapper(
	srcDir string,
	mode BatchMode,
	filter func(entry os.DirEntry) bool,
	entryFn func(entry os.DirEntry) error,
) error {

	entries, err := os.ReadDir(srcDir)

	if err != nil {
		return err
	}

	if !mode.Async {
		for _, entry := range entries {

			if !filter(entry) {
				fmt.Printf("skipped %s\n", entry.Name())
				continue
			}

			firstErr := entryFn(entry)
			if firstErr != nil {
				fmt.Printf("%v: %v\n", entry.Name(), firstErr)

				if mode.Strict {
					return firstErr
				}
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	var once sync.Once
	var mu sync.Mutex
	var selfErr error
	done := make(chan struct{})

	for _, entry := range entries {

		if !filter(entry) {
			fmt.Printf("skipped %s\n", entry.Name())
			continue
		}

		wg.Add(1)

		go func(entry os.DirEntry) {

			defer wg.Done()

			select {
			case <-done:
				return
			default:
			}

			firstErr := entryFn(entry)
			if firstErr != nil {
				fmt.Printf("%v: %v\n", entry.Name(), firstErr)

				mu.Lock()
				if selfErr == nil {
					selfErr = firstErr
				}
				mu.Unlock()

				if mode.Strict {
					once.Do(func() { close(done) })
				}
			}

		}(entry)
	}

	wg.Wait()
	return selfErr
}
