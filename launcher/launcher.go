package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Launcher struct {
	args    []string
	plugins []interface{}
	process *os.Process
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
	}
}

func (launcher *Launcher) Start() (err error) {
	args := launcher.args
	if args[0], err = exec.LookPath(args[0]); err == nil {
		stdin_r, stdin_w := launcher.openPipe()
		stdout_r, stdout_w := launcher.openPipe()
		stderr_r, stderr_w := launcher.openPipe()
		go connectPipes(os.Stdin, stdin_w, nil)
		go connectPipes(stdout_r, os.Stdout, func(contents []byte) {
			for _, plugin := range launcher.plugins {
				if handler, ok := plugin.(StdoutPlugin); ok {
					handler.ReadStdout(contents)
				}
			}
		})
		go connectPipes(stderr_r, os.Stderr, func(contents []byte) {
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

func (launcher *Launcher) InstallPlugin(plugin interface{}) {
	launcher.plugins = append(launcher.plugins, plugin)
}

func (launcher *Launcher) Cleanup(code int) {
	for _, plugin := range launcher.plugins {
		if handler, ok := plugin.(ShutdownPlugin); ok {
			handler.AtExit(code)
		}
	}
}

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
	launcher.Cleanup(exitCode)
	return nil
}

func (launcher *Launcher) openPipe() (*os.File, *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		launcher.Cleanup(-1)
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

func connectPipes(in *os.File, out *os.File, handler func([]byte)) {
	contents := make([]byte, 16)
	for {
		count, err := in.Read(contents)
		if count > 0 {
			out.Write(contents)
			if handler != nil {
				handler(contents)
			}
		}
		if err != nil {
			out.Close()
			break
		}
	}
}
