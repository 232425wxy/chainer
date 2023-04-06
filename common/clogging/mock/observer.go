package mock

import (
	"sync"

	"github.com/232425wxy/chainer/common/clogging"
	"go.uber.org/zap/zapcore"
)

type Observer struct {
	CheckStub        func(zapcore.Entry, *zapcore.CheckedEntry)
	checkMutex       sync.RWMutex
	checkArgsForCall []struct {
		arg1 zapcore.Entry
		arg2 *zapcore.CheckedEntry
	}
	WriteEntryStub        func(zapcore.Entry, []zapcore.Field)
	writeEntryMutex       sync.RWMutex
	writeEntryArgsForCall []struct {
		arg1 zapcore.Entry
		arg2 []zapcore.Field
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (o *Observer) Check(arg1 zapcore.Entry, arg2 *zapcore.CheckedEntry) {
	o.checkMutex.Lock()
	o.checkArgsForCall = append(o.checkArgsForCall, struct {
		arg1 zapcore.Entry
		arg2 *zapcore.CheckedEntry
	}{arg1: arg1, arg2: arg2})
	o.recordInvocation("Check", []interface{}{arg1, arg2})
	o.checkMutex.Unlock()
	if o.CheckStub != nil {
		o.CheckStub(arg1, arg2)
	}
}

func (o *Observer) CheckCallCount() int {
	o.checkMutex.RLock()
	defer o.checkMutex.RUnlock()
	return len(o.checkArgsForCall)
}

func (o *Observer) SetCheckStub(stub func(zapcore.Entry, *zapcore.CheckedEntry)) {
	o.checkMutex.Lock()
	o.CheckStub = stub
	o.checkMutex.Unlock()
}

func (o *Observer) CheckArgsForCall(i int) (zapcore.Entry, *zapcore.CheckedEntry) {
	o.checkMutex.RLock()
	defer o.checkMutex.RUnlock()
	arg1, arg2 := o.checkArgsForCall[i].arg1, o.checkArgsForCall[i].arg2
	return arg1, arg2
}

func (o *Observer) WriteEntry(arg1 zapcore.Entry, arg2 []zapcore.Field) {
	var arg2Copy []zapcore.Field
	if arg2 != nil {
		arg2Copy = make([]zapcore.Field, len(arg2))
		copy(arg2Copy, arg2)
	}
	o.writeEntryMutex.Lock()
	o.writeEntryArgsForCall = append(o.writeEntryArgsForCall, struct{arg1 zapcore.Entry; arg2 []zapcore.Field}{arg1: arg1, arg2: arg2Copy})
	o.recordInvocation("WriteEntry", []interface{}{arg1, arg2Copy})
	o.writeEntryMutex.Unlock()
	if o.WriteEntryStub != nil {
		o.WriteEntryStub(arg1, arg2)
	}
}

func (o *Observer) WriteEntryCallCount() int {
	o.writeEntryMutex.RLock()
	defer o.writeEntryMutex.RUnlock()
	return len(o.writeEntryArgsForCall)
}

func (o *Observer) SetWriteEntryStub(stub func(zapcore.Entry, []zapcore.Field)) {
	o.writeEntryMutex.Lock()
	o.WriteEntryStub = stub
	o.writeEntryMutex.Unlock()
}

func (o *Observer) WriteEntryArgsForCall(i int) (zapcore.Entry, []zapcore.Field) {
	o.writeEntryMutex.RLock()
	defer o.writeEntryMutex.RUnlock()
	arg1, arg2 := o.writeEntryArgsForCall[i].arg1, o.writeEntryArgsForCall[i].arg2
	return arg1, arg2
}

func (o *Observer) Invocations() map[string][][]interface{} {
	o.invocationsMutex.RLock()
	defer o.invocationsMutex.RUnlock()
	o.checkMutex.RLock()
	defer o.checkMutex.RUnlock()
	o.writeEntryMutex.RLock()
	defer o.writeEntryMutex.RUnlock()
	copiedInvocations := make(map[string][][]interface{})
	for key, value := range o.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (o *Observer) recordInvocation(method string, args []interface{}) {
	o.invocationsMutex.Lock()
	defer o.invocationsMutex.Unlock()
	if o.invocations == nil {
		o.invocations = make(map[string][][]interface{})
	}
	if o.invocations[method] == nil {
		o.invocations[method] = make([][]interface{}, 0)
	}
	o.invocations[method] = append(o.invocations[method], args)
}

var _ clogging.Observer = new(Observer)
