package quasar

import (
	"context"
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

	logName := ins.daemon.logName()
	or, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "cannot read cmd stdout")
	}
	go logging(or, logName)

	er, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "cannot read cmd stdout")
	}
	go logging(er, logName+".stderr")

	err = cmd.Start()
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

func (ins *rawCommandInstance) GetEnv(envname string) (string, error) {
	return "", errors.New("not implemented yet in raw command")
}

func (ins *rawCommandInstance) Close(closername string) error {
	return errors.New("not implemented yet in raw command")
}

type commandOpts struct {
	Launch string `yaml:"launch"`
}
