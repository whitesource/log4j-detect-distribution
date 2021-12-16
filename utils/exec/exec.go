package exec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"io"
	"os/exec"
)

// Commander is an interface for creating Commands
// it can be used for mocking out os command executions for testing
type Commander interface {
	// Command returns a Cmd instance which can be used to run a single command.
	Command(logger logr.Logger, executable string, args ...string) Command

	// CommandContext returns a Cmd instance which can be used to run a single command.
	//
	// The provided context is used to kill the process if the context becomes done
	// before the command completes on its own. For example, a timeout can be set in
	// the context.
	CommandContext(ctx context.Context, logger logr.Logger, executable string, args ...string) Command
}

// commander is an implementation of Commander, using exec.Cmd
type commander struct {
	logger logr.Logger
}

func New() Commander {
	return &commander{}
}

func (c commander) Command(logger logr.Logger, executable string, args ...string) Command {
	return newCommand(logger, exec.Command(executable, args...))
}

func (c commander) CommandContext(ctx context.Context, logger logr.Logger, executable string, args ...string) Command {
	return newCommand(logger, exec.CommandContext(ctx, executable, args...))
}

type Command interface {
	// Run is like Output, but without the stdout
	Run() error

	// Output executes the command, and returns the stdout
	Output() ([]byte, error)

	// SetDir sets the command's working directory
	SetDir(dir string)

	// SetEnv sets the environment variables for the command's execution environment
	SetEnv(env map[string]string)

	SetStdin(in io.Reader)
	SetStdout(out io.Writer)
	SetStderr(out io.Writer)
}

// command is an implementation of Command using exec.Cmd
type command struct {
	logger     logr.Logger
	cmd        *exec.Cmd
	stderr     bytes.Buffer
	stdout     bytes.Buffer
	stdoutFile string
	executed   bool
}

func newCommand(logger logr.Logger, cmd *exec.Cmd) *command {
	return &command{
		logger: logger.WithValues("cmd", cmd.Path, "args", cmd.Args),
		cmd:    cmd,
	}
}

func (c *command) Run() error {
	if c.executed {
		panic("Run(): command can't be executed more than once!")
	}
	c.cmd.Stderr = createMultiWriter(&c.stderr, c.cmd.Stderr)
	c.cmd.Stdout = createMultiWriter(&c.stdout, c.cmd.Stdout)
	c.logger.Info("beginning command execution")
	err := c.cmd.Run()
	c.logCommandCompletion(err, c.stdout.Bytes(), c.stderr.Bytes())
	return err
}

func (c *command) Output() ([]byte, error) {
	if c.executed {
		panic("Output(): command can't be executed more than once!")
	}
	c.cmd.Stderr = createMultiWriter(&c.stderr, c.cmd.Stderr)
	c.logger.Info("beginning command execution")
	output, err := c.cmd.Output()
	c.logCommandCompletion(err, output, c.stderr.Bytes())
	return output, err
}

func createMultiWriter(ws ...io.Writer) io.Writer {
	var writers []io.Writer
	for _, w := range ws {
		if w != nil {
			writers = append(writers, w)
		}
	}
	return io.MultiWriter(writers...)
}

func (c *command) SetStdoutFile(path string) {
	c.stdoutFile = path
}

func (c *command) SetDir(dir string) {
	c.cmd.Dir = dir
}

func (c *command) SetStdin(r io.Reader) {
	c.cmd.Stdin = r
}

func (c *command) SetStdout(w io.Writer) {
	c.cmd.Stdout = w
}

func (c *command) SetStderr(w io.Writer) {
	c.cmd.Stderr = w
}

func (c *command) SetEnv(env map[string]string) {
	if env == nil {
		return
	}
	for k, v := range env {
		e := fmt.Sprintf("%s=%s", k, v)
		c.cmd.Env = append(c.cmd.Env, e)
	}
}

// logCommandCompletion shows the result of the command execution
func (c *command) logCommandCompletion(err error, stdout []byte, stderr []byte) {
	log := c.logger.WithValues("stdout", string(stdout), "stderr", string(stderr))
	if err != nil {
		var er *exec.ExitError
		if errors.As(err, &er) {
			log = log.WithValues("code", er.ExitCode())
		}
		log.Error(err, "command execution failed")
	} else {
		log.Info("command execution success")
	}
}
