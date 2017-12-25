package gmap

import (
	"sync"
)

type StringInterfaceMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func NewStringInterfaceMap() *StringInterfaceMap {
	return &StringInterfaceMap{
		m: make(map[string]interface{}),
	}
}

// 哈希表克隆
func (this *StringInterfaceMap) Clone() *map[string]interface{} {
    m := make(map[string]interface{})
    this.RLock()
    for k, v := range this.m {
        m[k] = v
    }
    this.RUnlock()
    return &m
}

// 设置键值对
func (this *StringInterfaceMap) Set(key string, val interface{}) {
	this.Lock()
	this.m[key] = val
	this.Unlock()
}

// 批量设置键值对
func (this *StringInterfaceMap) BatchSet(m map[string]interface{}) {
    this.Lock()
    for k, v := range m {
        this.m[k] = v
    }
    this.Unlock()
}

// 获取键值
func (this *StringInterfaceMap) Get(key string) interface{} {
	this.RLock()
	val, _ := this.m[key]
	this.RUnlock()
	return val
}

// 删除键值对
func (this *StringInterfaceMap) Remove(key string) {
	this.Lock()
	delete(this.m, key)
	this.Unlock()
}

// 批量删除键值对
func (this *StringInterfaceMap) BatchRemove(keys []string) {
    this.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *StringInterfaceMap) GetAndRemove(key string) interface{} {
	this.Lock()
	val, exists := this.m[key]
	if exists {
		delete(this.m, key)
	}
	this.Unlock()
	return val
}

// 返回键列表
func (this *StringInterfaceMap) Keys() []string {
	this.RLock()
	keys := make([]string, 0)
	for key, _ := range this.m {
		keys = append(keys, key)
	}
    this.RUnlock()
	return keys
}

// 返回值列表(注意是随机排序)
func (this *StringInterfaceMap) Values() []interface{} {
	this.RLock()
	vals := make([]interface{}, 0)
	for _, val := range this.m {
		vals = append(vals, val)
	}
	this.RUnlock()
	return vals
}

// 是否存在某个键
func (this *StringInterfaceMap) Contains(key string) bool {
	this.RLock()
	_, exists := this.m[key]
	this.RUnlock()
	return exists
}

// 哈希表大小
func (this *StringInterfaceMap) Size() int {
	this.RLock()
	len := len(this.m)
	this.RUnlock()
	return len
}

// 哈希表是否为空
func (this *StringInterfaceMap) IsEmpty() bool {
	this.RLock()
	empty := (len(this.m) == 0)
	this.RUnlock()
	return empty
}

// 清空哈希表
func (this *StringInterfaceMap) Clear() {
    this.Lock()
    this.m = make(map[string]interface{})
    this.Unlock()
}