package quasar

import (
	"context"

	"github.com/pkg/errors"
)

type daemonType int

const (
	DockerDaemon daemonType = 1 + iota
	MySQLDaemon
	RedisDaemon
	RawCommandDaemon
)

var (
	daemonTypeMap = map[string]daemonType{
		"Docker":     DockerDaemon,
		"MySQL":      MySQLDaemon,
		"Redis":      RedisDaemon,
		"RawCommand": RawCommandDaemon,
	}
	unknownDaemonTypeError = errors.New("Unknown daemon type")
)

func (dt daemonType) ToInstance(d daemon) (instance, error) {
	switch dt {
	case RawCommandDaemon:
		return newRawCommandInstance(d), nil
	case MySQLDaemon:
		return newMySQLInstance(d), nil
	default:
		return nil, unknownDaemonTypeError
	}
}

type daemon struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`

	// for RawCommand
	Command commandOpts `yaml:"command"`
}

func (d daemon) ToInstance() (instance, error) {
	dt, ok := daemonTypeMap[d.Type]
	if !ok {
		return nil, unknownDaemonTypeError
	}
	return dt.ToInstance(d)
}

func (d daemon) logName() string {
	logName := d.Type
	if d.Name != "" {
		logName = d.Name
	}

	return logName
}

type instanceStatus int

const (
	instanceStatusUnknown instanceStatus = iota
	instanceStatusStopped
	instanceStatusLaunch
	instanceStatusRunning
	instanceStatusFailed
)

type instance interface {
	Run(ctx context.Context) error
	Stop() error
	Status() instanceStatus
	Wait() error
	GetEnv(string) (string, error)
	Close(string) error
}
