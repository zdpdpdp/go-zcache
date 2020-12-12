package zcache

import "time"

type Container struct {
	lru         *LRU
	singleFight *SingleFight
}

func New(maxSize int) *Container {
	return &Container{
		lru:         NewLRU(maxSize),
		singleFight: &SingleFight{},
	}
}

func (c *Container) Take(key Key, options ...EntryOptionsFunc) (Value, error) {
	value, ok := c.lru.Get(key)
	if ok {
		el := value.(*Entry)

		if el.needRefresh() {
			go c.singleFight.Do(key, el.options.loader)
		}

		el.accessTime = time.Now()

		if el.options.ValidCheck == nil {
			return el.Value, nil
		}

		if el.options.ValidCheck(el) {
			return el.Value, nil
		}
	}

	//key 不存在，或者校验不通过
	option := new(EntryOptions)
	for _, f := range options {
		f(option)
	}
	if option.loader == nil {
		return nil, ErrCacheMiss
	}
	loaderResp, err := c.singleFight.Do(key, option.loader)
	if err != nil {
		return nil, err
	}
	entry := Entry{
		options:    option,
		createTime: time.Now(),
		accessTime: time.Now(),
		Key:        key,
		Value:      loaderResp,
	}
	c.lru.Add(&entry)
	return loaderResp, nil
}

func (c *Container) Remove(key Key) bool {
	return c.lru.Remove(key)
}
