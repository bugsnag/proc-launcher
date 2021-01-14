package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Launcher struct {
	args    []string
	plugins []interface{}
	process *os.Process
	group   *sync.WaitGroup
}

type StdoutPlugin interface {
	ReadStdout(bytes []byte)
}

type StderrPlugin interface {
	ReadStderr(bytes []byte)
}

type ShutdownPlugin interface {
	AtExit(code int)
}

func New(args ...string) *Launcher {
	return &Launcher{
		args,
		make([]interface{}, 0),
		nil,
		&sync.WaitGroup{},
	}
}

// Run the launcher command, configuring pipe forwarding and plugins
func (launcher *Launcher) Start() (err error) {
	args := launcher.args
	if args[0], err = exec.LookPath(args[0]); err == nil {
		stdin_r, stdin_w := launcher.openPipe()
		stdout_r, stdout_w := launcher.openPipe()
		stderr_r, stderr_w := launcher.openPipe()
		launcher.group.Add(2)
		go connectPipes(nil, os.Stdin, stdin_w, nil)
		go connectPipes(launcher.group, stdout_r, os.Stdout, func(contents []byte) {
			for _, plugin := range launcher.plugins {
				if handler, ok := plugin.(StdoutPlugin); ok {
					handler.ReadStdout(contents)
				}
			}
		})
		go connectPipes(launcher.group, stderr_r, os.Stderr, func(contents []byte) {
			for _, plugin := range launcher.plugins {
				if handler, ok := plugin.(StderrPlugin); ok {
					handler.ReadStderr(contents)
				}
			}
		})

		notifications := make(chan os.Signal, 1)
		// will need to be in POSIX-specific file
		signal.Notify(notifications, syscall.SIGALRM, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go launcher.forwardSignals(notifications)

		var attrs os.ProcAttr
		attrs.Files = []*os.File{stdin_r, stdout_w, stderr_w}
		process, err := os.StartProcess(args[0], args, &attrs)
		if err == nil {
			launcher.process = process
			return nil
		}
	}
	return err
}

// Install a plugin which conforms to StdoutPlugin, StderrPlugin, and/or
// ShutdownPlugin.
func (launcher *Launcher) InstallPlugin(plugin interface{}) {
	launcher.plugins = append(launcher.plugins, plugin)
}

func (launcher *Launcher) cleanup(code int) {
	for _, plugin := range launcher.plugins {
		if handler, ok := plugin.(ShutdownPlugin); ok {
			handler.AtExit(code)
		}
	}
}

// Wait for the launched process to terminate. Should only be called after
// Start().
func (launcher *Launcher) Wait() error {
	if launcher.process == nil {
		return fmt.Errorf("process not yet started")
	}
	state, err := launcher.process.Wait()
	if err != nil {
		return fmt.Errorf("failed to await process: %v", err)
	}
	var exitCode int = 0
	if status, ok := state.Sys().(syscall.WaitStatus); ok {
		exitCode = status.ExitStatus()
	}
	launcher.group.Wait()
	launcher.cleanup(exitCode)
	return nil
}

func (launcher *Launcher) openPipe() (*os.File, *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		launcher.cleanup(-1)
		os.Stderr.WriteString(fmt.Sprintf("failed to open pipe: %v", err))
	}
	return r, w
}

func (launcher *Launcher) forwardSignals(notifications chan os.Signal) {
	for {
		signal := <-notifications
		if launcher.process != nil {
			launcher.process.Signal(signal)
		}
	}
}

func connectPipes(group *sync.WaitGroup, in *os.File, out *os.File, handler func([]byte)) {
	completedGroup := false
	defer func() {
		if group != nil && !completedGroup {
			group.Done()
		}
	}()
	contents := make([]byte, 512)
	shouldQuit := false
	for {
		now := time.Now()
		deadline := now.Add(time.Millisecond * 100)
		if shouldQuit {
			// grant extra time for final read since shutdown imminent
			deadline = deadline.Add(time.Millisecond * 400)
		}
		if err := in.SetReadDeadline(deadline); err == os.ErrNoDeadline && !completedGroup {
			if group != nil {
				completedGroup = true
				group.Done() // Read() may never return.
			}
		}
		count, err := in.Read(contents)
		for count > 0 {
			slice := contents[0:count]
			out.Write(slice)
			if handler != nil {
				handler(slice)
			}
			count, err = in.Read(contents)
		}
		if shouldQuit || (err != nil && !isDeadlineExceededErr(err)) {
			break
		}
		if count, err = in.Write([]byte{}); err != nil || count < 0 {
			shouldQuit = true
		}
	}
}
