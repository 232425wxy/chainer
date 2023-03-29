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

	NewGaugeStub          func(metrics.GaugeOpts) metrics.Gauge
	newGaugeMutex         sync.RWMutex
	newGaugeArgsForCall   []struct{ arg1 metrics.GaugeOpts }
	newGaugeReturns       struct{ result1 metrics.Gauge }
	newGaugeReturnsOnCall map[int]struct{ result1 metrics.Gauge }

	NewHistogramStub          func(metrics.HistogramOpts) metrics.Histogram
	newHistogramMutex         sync.RWMutex
	newHistogramArgsForCall   []struct{ arg1 metrics.HistogramOpts }
	newHistogramReturns       struct{ result1 metrics.Histogram }
	newHistogramReturnsOnCall map[int]struct{ result1 metrics.Histogram }

	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Provider) NewCounter(arg1 metrics.CounterOpts) metrics.Counter {
	fake.newCounterMutex.Lock()
	ret, specificReturn := fake.newCounterReturnsOnCall[len(fake.newCounterReturnsOnCall)]
	fake.newCounterArgsForCall = append(fake.newCounterArgsForCall, struct{ arg1 metrics.CounterOpts }{arg1})
	fake.recordInvocation("NewCounter", []interface{}{arg1})
	fake.newCounterMutex.Unlock()
	if fake.NewCounterStub != nil {
		return fake.NewCounterStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.newCounterReturns
	return fakeReturns.result1
}

// NewCounterCallCount 返回调用 NewCounter 方法的次数。
func (fake *Provider) NewCounterCallCount() int {
	fake.newCounterMutex.Lock()
	defer fake.newCounterMutex.Unlock()
	return len(fake.newCounterArgsForCall)
}

func (fake *Provider) NewGauge(arg1 metrics.GaugeOpts) metrics.Gauge {
	fake.newGaugeMutex.Lock()
	ret, specificReturn := fake.newGaugeReturnsOnCall[len(fake.newGaugeReturnsOnCall)]
	fake.newGaugeArgsForCall = append(fake.newGaugeArgsForCall, struct{arg1 metrics.GaugeOpts}{arg1})
	fake.recordInvocation("NewGauge", []interface{}{arg1})
}

func (fake *Provider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ metrics.Provider = new(Provider)
