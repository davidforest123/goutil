package gmap

import "goutil/container/gstring"

type (
	SortedMap struct {
		data      map[string]any
		queueKeys []string
	}
)

func NewSortedMap() *SortedMap {
	return &SortedMap{
		data:      map[string]any{},
		queueKeys: nil,
	}
}

func (sm *SortedMap) Has(key string) bool {
	_, exist := sm.data[key]
	return exist
}

func (sm *SortedMap) Delete(key string) {
	_, exist := sm.data[key]
	if !exist {
		return
	}
	gstring.RemoveByValue(sm.queueKeys, key)
	delete(sm.data, key)
}

func (sm *SortedMap) Upsert(key string, value any) {
	_, exist := sm.data[key]

	sm.data[key] = value
	if !exist {
		sm.queueKeys = append(sm.queueKeys, key)
	}
}

func (sm *SortedMap) Len() int {
	return len(sm.queueKeys)
}

func (sm *SortedMap) GetByKey(key string) any {
	return sm.data[key]
}

func (sm *SortedMap) GetByIndex(idx int) any {
	if idx < len(sm.queueKeys) {
		key := sm.queueKeys[idx]
		return sm.data[key]
	}
	return nil
}

func (sm *SortedMap) GetOldest() any {
	if len(sm.queueKeys) == 0 {
		return nil
	} else {
		return sm.data[sm.queueKeys[0]]
	}
}

func (sm *SortedMap) GetLatest() any {
	if len(sm.queueKeys) == 0 {
		return nil
	} else {
		return sm.data[sm.queueKeys[len(sm.queueKeys)-1]]
	}
}
