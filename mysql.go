package quasar

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync/atomic"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lestrrat/go-test-mysqld"
	"github.com/pkg/errors"
)

type mysqlInstance struct {
	daemon    daemon
	mysqld    *mysqltest.TestMysqld
	status    instanceStatus
	waitCh    chan struct{}
	dbnameCnt uint64
	logger    *bytes.Buffer
}

func newMySQLInstance(d daemon) *mysqlInstance {
	return &mysqlInstance{
		daemon:    d,
		status:    instanceStatusStopped,
		waitCh:    make(chan struct{}),
		dbnameCnt: 0,
		logger:    new(bytes.Buffer),
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
	go logging(ins.logger, logName)

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

func (ins *mysqlInstance) GetEnv(envname string) (string, error) {
	switch envname {
	case "dsn", "DSN":
		fmt.Fprintln(ins.logger, "received dsn request")
		return ins.createDatabase()
	default:
		return "", errors.New("undefined enviroment in MySQLDaemon")
	}
}

func (ins *mysqlInstance) createDatabase() (string, error) {
	md := ins.mysqld

	db, err := sql.Open("mysql", md.DSN())
	if err != nil {
		return "", errors.Wrap(err, "failed createDatabase")
	}
	defer db.Close()
	dbname := ins.generateDbname()
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		return "", errors.Wrap(err, "cannot create database in createDatabase")
	}

	fmt.Fprintf(ins.logger, "create database %s\n", dbname)

	dsn := md.DSN(mysqltest.WithDbname(dbname))
	return dsn, nil
}

func (ins *mysqlInstance) generateDbname() string {
	dbnameCnt := atomic.AddUint64(&ins.dbnameCnt, 1)
	return fmt.Sprintf("testdb%d", dbnameCnt)
}

func (ins *mysqlInstance) Close(dbname string) error {
	md := ins.mysqld
	db, err := sql.Open("mysql", md.DSN())
	if err != nil {
		return errors.Wrap(err, "failed connect db in Close")
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s", dbname))
	if err != nil {
		return errors.Wrap(err, "cannot drop database")
	}

	return nil
}
