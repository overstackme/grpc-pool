package ns

import (
	"context"

	"google.golang.org/grpc/resolver"
)

func NewBuilder(d Discovery) resolver.Builder {
	return &nsResolverBuilder{
		d: d,
	}
}

type nsResolverBuilder struct {
	d Discovery
}

func (nsb *nsResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &nsResolver{
		target: target,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
		d:      nsb.d,
	}
	go r.watcher()
	return r, nil
}

func (*nsResolverBuilder) Scheme() string {
	return "ns"
}
