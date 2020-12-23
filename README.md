# proc-launcher

An executable which launches a child process, awaiting its exit. Input/output
streams are forwarded to/from the child process, as well as system signals.

## Usage

```sh
proc-launcher [my_process] [my_process arguments]
```

## API Usage

The launcher can also be extended to handle terminated process state and output
stream contents

```go
type extension struct{}

// handle bytes added to stdout
func (ex extension) ReadStdout(bytes []byte) {
	fmt.Printf("you said: %s", string(bytes))
}

// handle bytes added to stderr
func (ex extension) ReadStderr(bytes []byte) {
	fmt.Printf("something bad happened: %s", string(bytes))
}

// handle child process termination
func (ex extension) AtExit(code int) {
	fmt.Printf("process terminated, code %d", code)
}

func main() {
  // create a new launcher using command arguments as the child process
	launcher := l.New(os.Args[1:]...)
  // install a custom extension. An extension can respond to any/all of the
  // above interfaces.
	launcher.InstallPlugin(extension{})
  // launch the child process
	if err := launcher.Start(); err != nil {
		fmt.Printf("failed to launch process: %v", err)
	}
  // wait until the process terminates to exit
	if err := launcher.Wait(); err != nil {
		fmt.Printf("failed to await process: %v", err)
	}
}
```

## Testing

```
cucumber
```
