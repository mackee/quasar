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
	default:
		return nil, unknownDaemonTypeError
	}
}

type daemon struct {
	Name    string      `yaml:"name"`
	Type    string      `yaml:"type"`
	Command commandOpts `yaml:"command"`
}

func (d daemon) ToInstance() (instance, error) {
	dt, ok := daemonTypeMap[d.Type]
	if !ok {
		return nil, unknownDaemonTypeError
	}
	return dt.ToInstance(d)
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
	PID() (int, error)
	Status() instanceStatus
	Wait() error
}
