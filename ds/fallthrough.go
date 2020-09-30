// +build !zookeeper

package ds

import "google.golang.org/grpc/resolver"

type Fallthrough struct{}

func NewDiscovery(config *Config) Discovery {
	return &Fallthrough{}
}

func (ds *Fallthrough) Start() error {
	return nil
}

func (ds *Fallthrough) Stop() {}

func (*Fallthrough) Addresses(target resolver.Target) ([]resolver.Address, error) {
	return []resolver.Address{
		{
			Addr: target.Endpoint,
		},
	}, nil
}
