package main

import (
	"fmt"
	"os"

	l "github.com/bugsnag/proc-launcher/launcher"
)

type extension struct{}

func (ex extension) ReadStdout(bytes []byte) {
	fmt.Printf("you said: %s", string(bytes))
}

func (ex extension) ReadStderr(bytes []byte) {
	fmt.Printf("something bad happened: %s", string(bytes))
}

func (ex extension) AtExit(code int) {
	fmt.Printf("process terminated, code %d", code)
}

func main() {
	launcher := l.New(os.Args[1:]...)
	launcher.InstallPlugin(extension{})
	if err := launcher.Start(); err != nil {
		fmt.Printf("failed to launch process: %v", err)
	}
	if err := launcher.Wait(); err != nil {
		fmt.Printf("failed to await process: %v", err)
	}
}
