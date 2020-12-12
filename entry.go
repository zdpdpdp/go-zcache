package zcache

import "time"

type Entry struct {
	options    *EntryOptions
	createTime time.Time
	accessTime time.Time
	Key        Key
	Value      Value
}

func (e *Entry) needRefresh() bool {
	return e.options.refresh > 0 && e.options.loader != nil && time.Now().After(e.accessTime.Add(e.options.refresh))
}

type (
	Key   interface{}
	Value interface{}
)
type Loader func(Key) (Value, error)

type ValidCondition func(e *Entry) bool

type EntryOptions struct {
	loader     Loader
	ExpireTime time.Time
	refresh    time.Duration

	ValidCheck ValidCondition //取数据时，会调用用户提供的自定义校验方法
}

type EntryOptionsFunc func(o *EntryOptions)

//WithExpire 写数据多久以后过期
func WithExpire(d time.Duration) EntryOptionsFunc {
	return func(o *EntryOptions) {
		o.ExpireTime = time.Now().Add(d)
	}
}

//WithRefresh 多久以后触发刷新调用loader写数据 lazy load
func WithRefresh(d time.Duration) EntryOptionsFunc {
	return func(o *EntryOptions) {
		o.refresh = time.Now().Add(d)
	}
}

//WithLoader loader加载数据，防缓存击穿， 多协程调用Get时只有一个协程会调用，其余协程返回旧数据，无旧数据时等待
func WithLoader(loader Loader) EntryOptionsFunc {
	return func(o *EntryOptions) {
		o.loader = loader
	}
}
