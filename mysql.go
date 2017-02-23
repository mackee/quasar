package quasar

import (
	"context"
	"os"

	"github.com/lestrrat/go-test-mysqld"
	"github.com/pkg/errors"
)

type mysqlInstance struct {
	daemon daemon
	mysqld *mysqltest.TestMysqld
	status instanceStatus
	waitCh chan struct{}
}

func newMySQLInstance(d daemon) *mysqlInstance {
	return &mysqlInstance{
		daemon: d,
		status: instanceStatusStopped,
		waitCh: make(chan struct{}),
	}
}

func (ins *mysqlInstance) Run(ctx context.Context) error {
	ins.status = instanceStatusLaunch
	mysqld, err := mysqltest.NewMysqld(nil)
	if err != nil {
		ins.status = instanceStatusFailed
		return errors.Wrap(err, "failed launch mysqld")
	}
	ins.mysqld = mysqld
	lf := mysqld.LogFile
	lr, err := os.Open(lf)
	if err != nil {
		ins.status = instanceStatusFailed
		ins.Stop()
		return errors.Wrap(err, "failed read logfile from mysqld")
	}
	logName := ins.daemon.logName()
	go logging(lr, logName)

	ins.status = instanceStatusRunning

	return nil
}

func (ins *mysqlInstance) Stop() error {
	ins.mysqld.Stop()
	ins.status = instanceStatusStopped

	ins.waitCh <- struct{}{}
	return nil
}

func (ins *mysqlInstance) Status() instanceStatus {
	return ins.status
}

func (ins *mysqlInstance) Wait() error {
	<-ins.waitCh
	return nil
}
