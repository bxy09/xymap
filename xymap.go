package xymap

import (
	"github.com/cheekybits/genny/generic"
)

type KeyType generic.Type
type ValueType generic.Type

type wrapperKeyTypeValueType struct {
	Key   KeyType
	Value ValueType
	Valid bool
}

type XYMapKeyTypeValueType struct {
	mapping map[KeyType]int
	storage []*wrapperKeyTypeValueType
	empty   int
}

func NewXYMapKeyTypeValueType() *XYMapKeyTypeValueType {
	return &XYMapKeyTypeValueType{
		mapping: map[KeyType]int{},
	}
}

func (xym *XYMapKeyTypeValueType) Set(key KeyType, value ValueType) (subs ValueType, existed bool) {
	if idx, ok := xym.mapping[key]; ok {
		slot := xym.storage[idx]
		if slot.Valid {
			subs = slot.Value
			existed = true
		}
		slot.Value = value
		return
	}
	xym.mapping[key] = len(xym.storage)
	xym.storage = append(xym.storage, &wrapperKeyTypeValueType{
		Key:   key,
		Value: value,
		Valid: true,
	})
	return
}

func (xym *XYMapKeyTypeValueType) Get(key KeyType) (value ValueType, existed bool) {
	if idx, ok := xym.mapping[key]; ok {
		slot := xym.storage[idx]
		if slot.Valid {
			value = slot.Value
			existed = true
		}
	}
	return
}

// Delete the key and the opposite value
func (xym *XYMapKeyTypeValueType) Delete(key KeyType) (value ValueType, existed bool) {
	if idx, ok := xym.mapping[key]; ok {
		slot := xym.storage[idx]
		if slot.Valid {
			value = slot.Value
			existed = true
			slot.Valid = false
			xym.empty++
			if xym.empty > 10 && float32(xym.empty) / float32(len(xym.storage)) > 0.8 {
				xym.Compress()
			}
		}
		return
	}
	return
}

// Compress the storage
func (xym *XYMapKeyTypeValueType) Compress() {
	wid := 0
	rid := len(xym.storage) - 1
	for ; rid >= wid; rid-- {
		slot := xym.storage[rid]
		if slot.Valid {
			for ; wid < rid && xym.storage[wid].Valid; wid++ {
			}
			delete(xym.mapping, xym.storage[wid].Key)
			xym.mapping[slot.Key] = wid
			xym.storage[wid] = slot
		} else {
			delete(xym.mapping, slot.Key)
		}
	}
	xym.storage = xym.storage[:wid + 1]
	xym.empty = 0
	return
}

func (xym *XYMapKeyTypeValueType) Iterate(callback func(key KeyType, value ValueType) bool) {
	length := len(xym.storage)
	for i := 0; i < length; i++ {
		slot := xym.storage[i]
		if slot.Valid {
			flag := callback(slot.Key, slot.Value)
			if flag {
				break
			}
		}
	}
}

func (xym *XYMapKeyTypeValueType) Length() int {
	return len(xym.storage) - xym.empty
}
