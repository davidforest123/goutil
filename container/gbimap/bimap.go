package gbimap

import (
	"errors"
	"reflect"
)

type (
	// BiMap is a bi-directional map
	BiMap struct {
		keys        map[any]any
		values      map[any]any
		orderedKeys []any
	}

	// Tuple is a key-value pair used to initialized a new BiMap with values
	Tuple struct {
		key, value any
	}
)

// Set sets a key-value pair on the map. Returns an error if key or value is duplicate
func (bm *BiMap) Set(key any, value any) error {
	if bm.keys[key] != nil {
		return errors.New("BiMap can\\'t have duplicate key")
	}

	if bm.values[value] != nil {
		return errors.New("BiMap can\\'t have duplicate value")
	}

	bm.keys[key] = value
	bm.values[value] = key
	bm.orderedKeys = append(bm.orderedKeys, key)

	return nil
}

// GetValByKey gets a value from a key
func (bm *BiMap) GetValByKey(key any) any {
	return bm.keys[key]
}

// GetKeyByVal gets a key from a value
func (bm *BiMap) GetKeyByVal(value any) any {
	return bm.values[value]
}

func (bm *BiMap) getIndexOfKey(key any) int {
	for i, k := range bm.orderedKeys {
		if key == k {
			return i
		}
	}

	return -1
}

func (bm *BiMap) deletePair(key any, value any) {
	delete(bm.keys, key)
	delete(bm.values, value)
	keyIndex := bm.getIndexOfKey(key)
	bm.orderedKeys = append(bm.orderedKeys[:keyIndex], bm.orderedKeys[:keyIndex+1])
}

// DelByVal deletes a key-value pair from a value. Returns an error if provided argument is not a value
func (bm *BiMap) DelByVal(value any) error {
	key := bm.GetKeyByVal(value)
	if key == nil {
		return errors.New("Key does not exist in BiMap")
	}

	bm.deletePair(key, value)
	return nil
}

// DelByKey deletes a key-value pair from a key. Returns an error if provided argument is not a key
func (bm *BiMap) DelByKey(key any) error {
	value := bm.GetValByKey(key)
	if value == nil {
		return errors.New("Value does not exist in BiMap")
	}

	bm.deletePair(key, value)
	return nil
}

// Size returns the size of the map
func (bm *BiMap) Size() int {
	return len(bm.keys)
}

// Left returns the "key: value" mapping of the BiMap
func (bm *BiMap) Left() map[any]any {
	return bm.keys
}

// Right returns the "value: key" mapping of the BiMap
func (bm *BiMap) Right() map[any]any {
	return bm.values
}

// Keys returns a slice with all the BiMap keys
func (bm *BiMap) Keys() []any {
	return bm.orderedKeys
}

// Vals returns a slice with all the BiMap values
func (bm *BiMap) Vals() []any {
	slice := make([]any, 0, len(bm.orderedKeys))
	for _, key := range bm.orderedKeys {
		slice = append(slice, bm.GetValByKey(key))
	}
	return slice
}

// IsEqual checks if a BiMap is equal to another
func (bm *BiMap) IsEqual(otherBm BiMap) bool {
	return reflect.DeepEqual(bm.values, otherBm.Right()) && reflect.DeepEqual(bm.keys, otherBm.Left())
}

// NewBiMap creates a new bi-directional map
func NewBiMap(initialValues ...Tuple) (*BiMap, error) {
	keys := make(map[any]any, len(initialValues))
	values := make(map[any]any, len(initialValues))
	orderedKeys := make([]any, 0, len(initialValues))

	for _, tuple := range initialValues {
		if keys[tuple.key] != nil {
			return nil, errors.New("Initial values contain duplicated keys")
		}

		if values[tuple.value] != nil {
			return nil, errors.New("Initial values contain duplicated values")
		}

		keys[tuple.key] = tuple.value
		values[tuple.value] = tuple.key
		orderedKeys = append(orderedKeys, tuple.key)
	}

	bm := &BiMap{
		keys:        keys,
		values:      values,
		orderedKeys: orderedKeys,
	}
	return bm, nil
}
