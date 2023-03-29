package metricsfakes

import (
	"sync"

	"github.com/232425wxy/chainer/common/metrics"
)

type Provider struct {
	NewCounterStub          func(metrics.CounterOpts) metrics.Counter
	newCounterMutex         sync.RWMutex
	newCounterArgsForCall   []struct{ arg1 metrics.CounterOpts }
	newCounterReturns       struct{ result1 metrics.Counter }
	newCounterReturnsOnCall map[int]struct{ result1 metrics.Counter }
}
