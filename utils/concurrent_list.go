package utils

import (
	"errors"
	"sync"
)

type List interface {
	Add(e interface{})
	Size() int
	Empty() bool
}

type ConcurrentList struct {
	lock    *sync.RWMutex
	entries []interface{}
}

func NewConcurrentList() *ConcurrentList {
	return NewConcurrentListSize(10)
}

func NewConcurrentListSize(size int) *ConcurrentList {
	list := new(ConcurrentList)

	list.lock = &sync.RWMutex{}
	list.entries = make([]interface{}, 0, size)

	return list
}

func (this *ConcurrentList) Add(e interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.entries = append(this.entries, e)
}

func (this *ConcurrentList) Size() int {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return len(this.entries)
}

func (this *ConcurrentList) Empty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return len(this.entries) == 0
}

func (this *ConcurrentList) RemoveFun(f func(int, interface{}) bool) []interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if len(this.entries) == 0 {
		return nil
	}

	indexs := make([]int, 0)
	result := make([]interface{}, 0)

	for index, e := range this.entries {

		if f(index, e) {
			indexs = append(indexs, index)
			result = append(result, e)
		}
	}

	this.entries = removeArray(this.entries, indexs...)

	return result
}

func (this *ConcurrentList) Remove(index int) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	size := len(this.entries)

	if size == 0 || index < 0 || index >= size {
		panic(errors.New("Outof range index"))
	}

	oldValue := this.entries[index]

	this.entries = removeArray(this.entries, index)

	return oldValue
}

func (this *ConcurrentList) Get(index int) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()

	size := len(this.entries)

	if size == 0 || index < 0 || index >= size {
		panic(errors.New("Outof range index"))
	}

	return this.entries[index]
}

func (this *ConcurrentList) ForEach(f func(int, interface{})) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	for index, e := range this.entries {
		f(index, e)
	}
}

func removeArray(entries []interface{}, indexs ...int) []interface{} {

	if entries == nil || len(entries) == 0 {
		return nil
	}

	if indexs == nil || len(indexs) == 0 {
		return entries
	}

	//remove all and initialize the entries
	if len(indexs) == len(entries) {
		return []interface{}{}
	}

	newEntries := make([]interface{}, len(entries)-len(indexs))

	marks := make(map[int]bool, 0)

	for _, j := range indexs {
		if j < 0 || j >= len(entries) {
			panic("index out of range")
		}

		marks[j] = true
	}

	i := 0

	for j, e := range entries {
		if !marks[j] {
			newEntries[i] = e
			i++
		}
	}

	return newEntries
}
