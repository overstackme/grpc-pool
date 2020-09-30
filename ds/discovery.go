package ds

import (
	"time"

	"github.com/overstackme/grpc-pool/ns"
)

const (
	defaultGroupName = "/grpc"
	defaultInterval  = 1 * time.Second
)

type Config struct {
	Endpoints []string
	Timeout   time.Duration
	GroupName string
}

type Service struct {
	Name string
	Addr string
}

type Discovery interface {
	ns.Discovery
	Watcher() error
	Registry([]Service) error
	Stop()
}
