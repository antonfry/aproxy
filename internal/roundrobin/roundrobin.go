package roundrobin

import (
	"aproxy/internal/backend"
	"aproxy/internal/targetgroup"
	"errors"
	"sync/atomic"
)

var ErrNoAliveUDPBackends = errors.New("no alive UDPBackends")

type Pool struct {
	Current     atomic.Int32
	TargetGroup *targetgroup.TargetGroup
}

func New(tg *targetgroup.TargetGroup) *Pool {
	return &Pool{
		Current:     atomic.Int32{},
		TargetGroup: tg,
	}
}

func (p *Pool) Next() (*backend.UDPBackend, error) {
	if len(p.TargetGroup.Backends) == 0 {
		return nil, ErrNoAliveUDPBackends
	}
	c := p.Current.Load()
	if int(c) >= len(p.TargetGroup.Backends) {
		c = 0
		p.Current.Store(0)
	}
	for n, b := range p.TargetGroup.Backends[c:] {
		if b.IsAlive() {
			p.Current.Store(c + 1 + int32(n))
			return b, nil
		}
	}
	for n, b := range p.TargetGroup.Backends[:c] {
		if b.IsAlive() {
			p.Current.Store(c + 1 + int32(n))
			return b, nil
		}
	}
	return nil, ErrNoAliveUDPBackends
}
