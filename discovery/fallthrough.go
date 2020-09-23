package discovery

import "google.golang.org/grpc/resolver"

type Fallthrough struct{}

func NewFallthrough() *Fallthrough {
	return &Fallthrough{}
}

func (*Fallthrough) Addresses(target resolver.Target) ([]resolver.Address, error) {
	return []resolver.Address{
		{
			Addr: target.Endpoint,
		},
	}, nil
}
