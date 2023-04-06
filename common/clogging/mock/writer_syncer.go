package mock

import "sync"

type WriteSyncer struct {
	SyncStub        func() error
	syncMutex       sync.RWMutex
	syncArgsForCall []struct{}
	syncReturns     struct {
		result1 error
	}
	syncReturnsOnCall map[int]struct {
		result1 error
	}
	WriteStub func([]byte) (int, error)
	writeMutex sync.RWMutex
	writeArgsForCall []struct {
		arg1 []byte
	}
	writeReturns struct {
		result1 int
		result2 error
	}
	writeReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (ws *WriteSyncer) Sync() error {
	ws.syncMutex.Lock()
	ret, specificReturn := ws.syncReturnsOnCall[len(ws.syncArgsForCall)] // 最后一次调用的返回值。
	ws.syncArgsForCall = append(ws.syncArgsForCall, struct{}{})
	ws.recordInvocation("Sync", []interface{}{})
	ws.syncMutex.Unlock()
	if ws.SyncStub != nil {
		return ws.SyncStub()
	}
	if specificReturn {
		return ret.result1
	}
	return ws.syncReturns.result1
}

func (ws *WriteSyncer) SyncCallCount() int {
	ws.syncMutex.RLock()
	defer ws.syncMutex.RUnlock()
	return len(ws.syncArgsForCall)
}

func (ws *WriteSyncer) SetSyncStub(stub func() error) {
	ws.syncMutex.Lock()
	ws.SyncStub = stub
	ws.syncMutex.Unlock()
}

func (ws *WriteSyncer) SetSyncReturns(result1 error) {
	ws.syncMutex.Lock()
	defer ws.syncMutex.Unlock()
	ws.SyncStub = nil
	ws.syncReturns = struct {result1 error}{result1: result1}
}

func (ws *WriteSyncer) SyncReturnsOnCall(i int, result1 error) {
	ws.syncMutex.Lock()
	defer ws.syncMutex.Unlock()
	ws.SyncStub = nil
	if ws.syncReturnsOnCall == nil {
		ws.syncReturnsOnCall = make(map[int]struct{result1 error})
	}
	ws.syncReturnsOnCall[i] = struct{result1 error}{result1: result1}
}

func (ws *WriteSyncer) Write(arg1 []byte) (int, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	ws.writeMutex.Lock()
	ret, specificReturn := ws.writeReturnsOnCall[len(ws.writeArgsForCall)]
	ws.writeArgsForCall = append(ws.writeArgsForCall, struct{arg1 []byte}{arg1: arg1})
	ws.recordInvocation("Write", []interface{}{arg1Copy})
	ws.writeMutex.Unlock()
	if ws.WriteStub != nil {
		return ws.WriteStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return ws.writeReturns.result1, ws.writeReturns.result2
}

func (ws *WriteSyncer) WriteCallCount() int {
	ws.writeMutex.RLock()
	defer ws.writeMutex.RUnlock()
	return len(ws.writeArgsForCall)
}

func (ws *WriteSyncer) SetWriteStub(stub func([]byte) (int, error)) {
	ws.writeMutex.Lock()
	ws.WriteStub = stub
	ws.writeMutex.Unlock()
}

func (ws *WriteSyncer) WriteArgsForCall(i int) []byte {
	ws.writeMutex.RLock()
	defer ws.writeMutex.RUnlock()
	return ws.writeArgsForCall[i].arg1
}

func (ws *WriteSyncer) SetWriteReturns(result1 int, result2 error) {
	ws.writeMutex.Lock()
	defer ws.writeMutex.Unlock()
	ws.WriteStub = nil
	ws.writeReturns = struct{result1 int; result2 error}{result1: result1, result2: result2}
}

func (ws *WriteSyncer) SetWriteReturnsOnCall(i int, result1 int, result2 error) {
	ws.writeMutex.Lock()
	defer ws.writeMutex.Unlock()
	ws.WriteStub = nil
	if ws.writeReturnsOnCall == nil {
		ws.writeReturnsOnCall = make(map[int]struct{result1 int; result2 error})
	}
	ws.writeReturnsOnCall[i] = struct{result1 int; result2 error}{result1: result1, result2: result2}
}

func (ws *WriteSyncer) Invocations() map[string][][]interface{} {
	ws.syncMutex.RLock()
	defer ws.syncMutex.RUnlock()
	ws.writeMutex.RLock()
	defer ws.writeMutex.RUnlock()
	ws.invocationsMutex.RLock()
	defer ws.invocationsMutex.RUnlock()
	copiedInvocations := make(map[string][][]interface{})
	for key, value := range ws.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (ws *WriteSyncer) recordInvocation(method string, args []interface{}) {
	ws.invocationsMutex.Lock()
	defer ws.invocationsMutex.Unlock()
	if ws.invocations == nil {
		ws.invocations = make(map[string][][]interface{})
	}
	if ws.invocations[method] == nil {
		ws.invocations[method] = make([][]interface{}, 0)
	}
	ws.invocations[method] = append(ws.invocations[method], args)
}
