package ns

import (
	"context"
	"time"

	"google.golang.org/grpc/resolver"
)

const (
	nsDefaultSyncInterval = 2 * time.Second
)

type Discovery interface {
	Addresses(resolver.Target) ([]resolver.Address, error)
}

type nsResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
	d      Discovery
}

func (r *nsResolver) watcher() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}

		state, err := r.lookup()
		if err != nil {
			r.cc.ReportError(err)
		} else {
			r.cc.UpdateState(*state)
		}

		t := time.NewTimer(nsDefaultSyncInterval)
		select {
		case <-t.C:
		case <-r.ctx.Done():
			t.Stop()
			return
		}
	}
}

func (r *nsResolver) lookup() (*resolver.State, error) {
	resolverAddress, err := r.d.Addresses(r.target)
	if err != nil {
		return nil, err
	}
	state := resolver.State{Addresses: resolverAddress}
	return &state, nil
}

func (*nsResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (r *nsResolver) Close() {
	r.cancel()
}
