package gmap

import "sync"

func SyncMapLen(m *sync.Map) int {
	var len int
	m.Range(func(k, v interface{}) bool {
		len++
		return true
	})
	return len
}
