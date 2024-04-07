package gbuiltin

import "sync/atomic"

type (
	AtomicBool int32
)

func NewAtomicBool(b bool) *AtomicBool {
	v := int32(0)
	if b {
		v = 1
	}
	result := new(int32)
	atomic.StoreInt32(result, v)
	return (*AtomicBool)(result)
}

func (ab *AtomicBool) Set(b bool) {
	v := int32(0)
	if b {
		v = 1
	}
	atomic.StoreInt32((*int32)(ab), v)
}

func (ab *AtomicBool) Get() bool {
	return atomic.LoadInt32((*int32)(ab)) != 0
}
