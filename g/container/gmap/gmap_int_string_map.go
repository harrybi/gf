// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gmap

import (
    "gitee.com/johng/gf/g/container/internal/rwmutex"
)

type IntStringMap struct {
	mu *rwmutex.RWMutex
	m  map[int]string
}

func NewIntStringMap(safe...bool) *IntStringMap {
	return &IntStringMap{
        m  : make(map[int]string),
        mu : rwmutex.New(safe...),
    }
}

// 给定回调函数对原始内容进行遍历，回调函数返回true表示继续遍历，否则停止遍历
func (this *IntStringMap) Iterator(f func (k int, v string) bool) {
    this.mu.RLock()
    defer this.mu.RUnlock()
    for k, v := range this.m {
        if !f(k, v) {
            break
        }
    }
}

// 哈希表克隆
func (this *IntStringMap) Clone() map[int]string {
	m := make(map[int]string)
	this.mu.RLock()
	for k, v := range this.m {
		m[k] = v
	}
    this.mu.RUnlock()
	return m
}

// 设置键值对
func (this *IntStringMap) Set(key int, val string) {
	this.mu.Lock()
	this.m[key] = val
	this.mu.Unlock()
}

// 批量设置键值对
func (this *IntStringMap) BatchSet(m map[int]string) {
	this.mu.Lock()
	for k, v := range m {
		this.m[k] = v
	}
	this.mu.Unlock()
}

// 获取键值
func (this *IntStringMap) Get(key int) string {
	this.mu.RLock()
	val, _ := this.m[key]
	this.mu.RUnlock()
	return val
}

// 设置kv缓存键值对，内部会对键名的存在性使用写锁进行二次检索确认，如果存在则不再写入；返回键名对应的键值。
// 在高并发下有用，防止数据写入的并发逻辑错误。
func (this *IntStringMap) doSetWithLockCheck(key int, value string) string {
    this.mu.Lock()
    if v, ok := this.m[key]; ok {
        this.mu.Unlock()
        return v
    }
    this.m[key] = value
    this.mu.Unlock()
    return value
}

// 当键名存在时返回其键值，否则写入指定的键值
func (this *IntStringMap) GetOrSet(key int, value string) string {
    this.mu.RLock()
    v, ok := this.m[key]
    this.mu.RUnlock()
    if !ok {
        return this.doSetWithLockCheck(key, value)
    } else {
        return v
    }
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func (this *IntStringMap) GetOrSetFunc(key int, f func() string) string {
    this.mu.RLock()
    v, ok := this.m[key]
    this.mu.RUnlock()
    if !ok {
        return this.doSetWithLockCheck(key, f())
    } else {
        return v
    }
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func (this *IntStringMap) GetOrSetFuncLock(key int, f func() string) string {
    this.mu.RLock()
    val, ok := this.m[key]
    this.mu.RUnlock()
    if !ok {
        this.mu.Lock()
        defer this.mu.Unlock()
        if v, ok := this.m[key]; ok {
            this.mu.Unlock()
            return v
        }
        val         = f()
        this.m[key] = val
        return val
    } else {
        return val
    }
}

// 当键名不存在时写入，并返回true；否则返回false。
func (this *IntStringMap) SetIfNotExist(key int, value string) bool {
    if !this.Contains(key) {
        this.doSetWithLockCheck(key, value)
        return true
    }
    return false
}

// 批量删除键值对
func (this *IntStringMap) BatchRemove(keys []int) {
    this.mu.Lock()
    for _, key := range keys {
        delete(this.m, key)
    }
    this.mu.Unlock()
}

// 返回对应的键值，并删除该键值
func (this *IntStringMap) Remove(key int) string {
    this.mu.Lock()
    val, exists := this.m[key]
    if exists {
        delete(this.m, key)
    }
    this.mu.Unlock()
    return val
}

// 返回键列表
func (this *IntStringMap) Keys() []int {
    this.mu.RLock()
    keys := make([]int, 0)
    for key, _ := range this.m {
        keys = append(keys, key)
    }
    this.mu.RUnlock()
    return keys
}

// 返回值列表(注意是随机排序)
func (this *IntStringMap) Values() []string {
    this.mu.RLock()
    vals := make([]string, 0)
    for _, val := range this.m {
        vals = append(vals, val)
    }
    this.mu.RUnlock()
    return vals
}

// 是否存在某个键
func (this *IntStringMap) Contains(key int) bool {
    this.mu.RLock()
    _, exists := this.m[key]
    this.mu.RUnlock()
    return exists
}

// 哈希表大小
func (this *IntStringMap) Size() int {
    this.mu.RLock()
    length := len(this.m)
    this.mu.RUnlock()
    return length
}

// 哈希表是否为空
func (this *IntStringMap) IsEmpty() bool {
    this.mu.RLock()
    empty := len(this.m) == 0
    this.mu.RUnlock()
    return empty
}

// 清空哈希表
func (this *IntStringMap) Clear() {
    this.mu.Lock()
    this.m = make(map[int]string)
    this.mu.Unlock()
}

// 并发安全锁操作，使用自定义方法执行加锁修改操作
func (this *IntStringMap) LockFunc(f func(m map[int]string)) {
    this.mu.Lock(true)
    defer this.mu.Unlock(true)
    f(this.m)
}

// 并发安全锁操作，使用自定义方法执行加锁读取操作
func (this *IntStringMap) RLockFunc(f func(m map[int]string)) {
    this.mu.RLock(true)
    defer this.mu.RUnlock(true)
    f(this.m)
}
