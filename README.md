# go-zcache

go本地缓存组件

## 功能

1. lru 缓存容量配置
2. 缓存防击穿
3. 缓存过期策略
    - expireAfterWrite 写数据多久后过期
    - refreshAfterWrite 数据多久后refresh,过期后第一个请求会触发异步更新缓存。但会返回旧的值，直到缓存被更新
    - expireAfterAccess 多久不访问就会过期
    - expireCondition  自定义过期判断方法
4. 命中统计

## 使用方法

```golang



```