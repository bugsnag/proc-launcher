package main

import (
	"fmt"
	"github.com/bugsnag/proc-launcher/launcher"
	"os"
)

func main() {
	launcher := launcher.New(os.Args[1:]...)
	if err := launcher.Start(); err != nil {
		fatalf("failed to launch process: %v", err)
	}
	if err := launcher.Wait(); err != nil {
		fatalf("failed to await process: %v", err)
	}
}

func fatalf(format string, args ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(format, args...))
	os.Exit(1)
}
