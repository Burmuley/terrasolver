package main

import (
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Queue interface {
	Next() *ExecModule
}

type Module interface {
	Exec(string, ...string) ([]byte, error)
}

type ExecModule struct {
	Path      string
	Processed bool
	Stdout    []byte
	Stderr    []byte
	ExitCode  int
}

func (m *ExecModule) Exec(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = m.Path
	pt, startErr := pty.Start(cmd)
	defer func() { _ = pt.Close() }()

	if startErr != nil {
		return startErr
	}

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, pt); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and back to Stdout
	go func() { _, _ = io.Copy(pt, os.Stdin) }()
	go func() { _, _ = io.Copy(os.Stdout, pt) }()

	return cmd.Wait()
}

type ExecQueue struct {
	modules []*ExecModule
	current int // index in the modules slice of the currently running module
}

func (e *ExecQueue) Next() *ExecModule {
	if e.current >= len(e.modules)-1 {
		return nil
	}

	e.current++
	return e.modules[e.current]
}

func NewExecQueue(modules []string) *ExecQueue {
	q := &ExecQueue{
		modules: make([]*ExecModule, 0, len(modules)),
	}

	for _, m := range modules {
		exm := &ExecModule{
			Path:      m,
			Processed: false,
		}

		q.modules = append(q.modules, exm)
		q.current = -1
	}

	return q
}
