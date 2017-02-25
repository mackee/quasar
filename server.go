package quasar

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/pkg/errors"
)

var (
	DaemonNameIsNotFoundError = errors.New("daemon name is not found.")
)

type Server struct {
	instances map[string]instance
}

type RPCArgs struct {
	Name, Envname string
}

func (s *Server) GetEnv(args *RPCArgs, resp *string) error {
	ins, ok := s.instances[args.Name]
	if !ok {
		return DaemonNameIsNotFoundError
	}

	r, err := ins.GetEnv(args.Envname)
	if err != nil {
		return errors.Wrapf(err, "failed call %s.GetEnv", args.Name)
	}

	*resp = r

	return nil
}

func (s *Server) Close(args *RPCArgs, resp *string) error {
	ins, ok := s.instances[args.Name]
	if !ok {
		return DaemonNameIsNotFoundError
	}

	err := ins.Close(args.Envname)
	if err != nil {
		return errors.Wrapf(err, "failed call %s.Close", args.Name)
	}

	return nil
}

func Serve(c config, inss map[string]instance) error {
	server := &Server{
		instances: inss,
	}
	rpc.Register(server)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", c.Address())
	if err != nil {
		return errors.Wrap(err, "cannot open server port")
	}

	return http.Serve(l, nil)
}
