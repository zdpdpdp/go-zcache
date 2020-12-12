package zcache

import "sync"

type call struct {
	wg  sync.WaitGroup
	val Value
	err error
}

type SingleFight struct {
	sync.Mutex
	calls map[Key]*call
}

func (g *SingleFight) Do(key Key, loader Loader) (Value, error) {
	g.Lock()
	if g.calls == nil {
		g.calls = make(map[Key]*call)
	}
	if c, ok := g.calls[key]; ok {
		//有其他协程正在call
		g.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.Unlock()

	c.val, c.err = loader(key)
	c.wg.Done()

	g.Lock()
	delete(g.calls, key)
	g.Unlock()
	return c.val, c.err
}
