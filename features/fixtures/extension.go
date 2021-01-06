package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	l "github.com/bugsnag/proc-launcher/launcher"
)

type extension struct {
	buffer *bufio.Writer
	group  *sync.WaitGroup
}

func (ex extension) ReadStdout(bytes []byte) {
	ex.write("you said: %s\n", string(bytes))
}

func (ex extension) ReadStderr(bytes []byte) {
	ex.write("something bad happened: %s\n", string(bytes))
}

func (ex extension) AtExit(code int) {
	ex.write("process terminated, code %d\n", code)
	ex.group.Done()
}

func (ex extension) write(format string, args ...interface{}) {
	defer ex.buffer.Flush()

	ex.buffer.WriteString(fmt.Sprintf(format, args...))
}

func New() extension {
	ex := extension{
		bufio.NewWriter(os.Stdout),
		&sync.WaitGroup{},
	}

	ex.group.Add(1)
	return ex
}

func main() {
	launcher := l.New(os.Args[1:]...)
	ex := New()
	launcher.InstallPlugin(ex)
	if err := launcher.Start(); err != nil {
		fmt.Printf("failed to launch process: %v\n", err)
	}
	if err := launcher.Wait(); err != nil {
		fmt.Printf("failed to await process: %v\n", err)
	}
	ex.group.Wait()
}
