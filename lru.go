package zcache

import (
	"container/list"
	"sync"
)

type LRU struct {
	len  int
	cap  int
	ls   list.List
	data sync.Map
	sync.RWMutex
}

func NewLRU(cap int) *LRU {
	l := &LRU{
		cap: cap,
	}
	l.ls.Init()
	return l
}

func (l *LRU) Add(e *Entry) {
	l.Lock()
	defer l.Unlock()

	ele := l.ls.PushFront(e)
	l.len++
	if l.len > l.cap {
		back := l.ls.Back()
		key := back.Value.(*Entry).Key
		l.data.Delete(key)
		l.ls.Remove(back)
		l.len--
	}
	l.data.Store(e.Key, ele)
}

func (l *LRU) Get(key Key) (Value, bool) {
	value, ok := l.data.Load(key)

	ele := value.(*list.Element)

	if !ok {
		return nil, false
	}
	l.Lock()
	defer l.Unlock()

	l.ls.MoveToFront(ele)
	return ele.Value, true
}

func (l *LRU) Remove(key Key) bool {
	l.Lock()
	defer l.Unlock()

	el, ok := l.data.Load(key)
	if !ok {
		return false
	}
	l.data.Delete(key)
	l.ls.Remove(el.(*list.Element))
	l.len--
	return true
}
