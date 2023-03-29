package metricsfakes

import (
	"sync"

	"github.com/232425wxy/chainer/common/metrics"
)

type Counter struct {
	AddStub        func(float64)
	addMutex       sync.RWMutex
	addArgsForCall []struct{ arg1 float64 }

	WithStub          func(...string) metrics.Counter
	withMutex         sync.RWMutex
	withArgsForCall   []struct{ arg1 []string }
	withReturns       struct{ result1 metrics.Counter }
	withReturnsOnCall map[int]struct{ result1 metrics.Counter }

	invocations       map[string][][]interface{}
	invocationsMutex  sync.RWMutex
}

func (fake *Counter) Add(arg1 float64) {
	fake.addMutex.Lock()
	fake.addArgsForCall = append(fake.addArgsForCall, struct{arg1 float64}{arg1})
	fake.recordInvocation("Add", []interface{}{arg1})
	fake.addMutex.Unlock()
	if fake.AddStub != nil {
		fake.AddStub(arg1)
	}
}

// AddCallCount 返回调用 Add 方法的次数。
func (fake *Counter) AddCallCount() int {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	return len(fake.addArgsForCall)
}

// AddCalls 设置 AddStub，func(float64)。
func (fake *Counter) AddCalls(stub func(float64)) {
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.AddStub = stub
}

// AddArgsForCall 返回第 i+1 次调用 Add 方法传入的参数，float64。
func (fake *Counter) AddArgsForCall(i int) float64 {
	fake.addMutex.RLock()
	defer fake.addMutex.RUnlock()
	argsForCall := fake.addArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Counter) With(arg1 ...string) metrics.Counter {
	fake.withMutex.Lock()
	ret, specifiedReturn := fake.withReturnsOnCall[len(fake.withArgsForCall)]
	fake.withArgsForCall = append(fake.withArgsForCall, struct{arg1 []string}{arg1})
	fake.recordInvocation("With", []interface{}{arg1})
	fake.withMutex.Unlock()
	if fake.WithStub != nil {
		return fake.WithStub(arg1...)
	}
	if specifiedReturn {
		return ret.result1
	}
	fakeReturns := fake.withReturns
	return fakeReturns.result1
}

// WithCallCount 返回 With 方法被调用的次数。
func (fake *Counter) WithCallCount() int {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	return len(fake.withArgsForCall)
}

// WithCalls 设置 WithStub，func(...string) metrics.Counter。
func (fake *Counter) WithCalls(stub func(...string) metrics.Counter) {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	fake.WithStub = stub
}

// WithArgsForCall 返回第 i+1 次调用 With 方法传入的参数，[]string。
func (fake *Counter) WithArgsForCall(i int) []string {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	argsForCall := fake.withArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Counter) WithReturns(result1 metrics.Counter) {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	fake.WithStub = nil
	fake.withReturns = struct{result1 metrics.Counter}{result1}
}

func (fake *Counter) WithReturnsOnCall(i int, result1 metrics.Counter) {
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	fake.WithStub = nil
	if fake.withReturnsOnCall == nil {
		fake.withReturnsOnCall = make(map[int]struct{result1 metrics.Counter})
	}
	fake.withReturnsOnCall[i] = struct{result1 metrics.Counter}{result1}
}

func (fake *Counter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	fake.addMutex.Lock()
	defer fake.addMutex.Unlock()
	fake.withMutex.Lock()
	defer fake.withMutex.Unlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Counter) recordInvocation(key string, args []interface{}) {
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

var _ metrics.Counter = new(Counter)