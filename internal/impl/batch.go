package impl

import (
	"fmt"
	"github.com/XV-521/fileops/v2/core"
	"github.com/XV-521/fileops/v2/core/out"
	"os"
	"sync"
)

type BatchMode struct {
	Sem    int
	Strict bool
	Rec    bool
}

func DoBatchWrap(
	srcDir string,
	mode BatchMode,
	filter func(ei core.EntryInfo) bool,
	entryFn func(ei core.EntryInfo) error,
) error {

	if mode.Sem == 0 {
		return fmt.Errorf("mode.Sem is 0")
	}

	var wg sync.WaitGroup
	var once sync.Once
	var mu sync.Mutex
	var selfErr error

	done := make(chan struct{})
	sem := make(chan struct{}, mode.Sem)

	var walk func(dir string)

	walk = func(dir string) {
		select {
		case <-done:
			return
		default:
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			mu.Lock()
			if selfErr == nil {
				selfErr = err
			}
			mu.Unlock()

			if mode.Strict {
				once.Do(func() { close(done) })
			}
			return
		}

		for _, entry := range entries {

			ei := core.EntryInfo{
				Dir:  dir,
				View: entry,
			}

			if ei.IsDir() && mode.Rec {
				walk(ei.Path())
				continue
			}

			if !filter(ei) {
				out.Skip.Printf("skipped %v\n", ei.Path())
				continue
			}

			wg.Add(1)

			go func(ei core.EntryInfo) {

				defer wg.Done()

				select {
				case <-done:
					return
				case sem <- struct{}{}:
				}

				defer func() { <-sem }()

				firstErr := entryFn(ei)

				if firstErr != nil {
					if !mode.Strict {
						out.Warn.Printf("failed %v: %v\n", ei.Path(), firstErr)
					}

					mu.Lock()
					if selfErr == nil {
						selfErr = firstErr
					}
					mu.Unlock()

					if mode.Strict {
						once.Do(func() { close(done) })
					}
				}

			}(ei)
		}
	}

	walk(srcDir)

	wg.Wait()

	return selfErr
}
