package quasar

import (
	"context"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

type rawCommandInstance struct {
	daemon daemon
	cmd    *exec.Cmd
}

func newRawCommandInstance(d daemon) *rawCommandInstance {
	return &rawCommandInstance{
		daemon: d,
	}
}

func (ins *rawCommandInstance) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", ins.daemon.Command.Launch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "fail start instance")
	}

	ins.cmd = cmd

	return nil
}

func (ins *rawCommandInstance) Stop() error {
	cmd := ins.cmd
	if cmd == nil {
		return errors.New("not start instance")
	}

	err := cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return errors.Wrap(err, "fail send stopping singal")
	}

	return nil
}

func (ins *rawCommandInstance) PID() (int, error) {
	if ins.cmd == nil {
		return 0, errors.New("not start instance")
	}
	if ins.cmd.ProcessState.Exited() {
		return 0, errors.New("process is exited")
	}
	return ins.cmd.Process.Pid, nil
}

func (ins *rawCommandInstance) Status() instanceStatus {
	if ins.cmd == nil {
		return instanceStatusStopped
	}
	if ins.cmd.ProcessState != nil && ins.cmd.ProcessState.Exited() {
		return instanceStatusFailed
	}

	return instanceStatusRunning
}

func (ins *rawCommandInstance) Wait() error {
	return ins.cmd.Wait()
}

type commandOpts struct {
	Launch string `yaml:"launch"`
}
