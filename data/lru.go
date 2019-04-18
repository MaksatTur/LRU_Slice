package data

import (
	"sync"
)

// ToDo ...
type ToDo struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
}

// LRU ...
type LRU struct {
	sync.Mutex
	Capacity  int
	ToDoItems map[string]ToDo `json:"ToDoItems"`
	Queue     []string
}

//NewLRU ...
func NewLRU(capacity int) *LRU {
	return &LRU{
		Capacity:  capacity,
		ToDoItems: make(map[string]ToDo),
		Queue:     make([]string, 0, capacity),
	}
}

func getSliceIndex(key string, sl ...string) int {
	for i := 0; i < len(sl); i++ {
		if sl[i] == key {
			return i
		}
	}
	return -1
}

func (lru *LRU) setElementOnHead(key string, index int) {
	lru.Lock()
	defer lru.Unlock()
	tmp := lru.Queue[index]
	lru.Queue[index] = lru.Queue[0]
	lru.Queue[0] = tmp

}

func (lru *LRU) purge() {
	if len(lru.Queue) > 0 {
		delete(lru.ToDoItems, lru.Queue[len(lru.Queue)-1])
		lru.Queue = lru.Queue[:len(lru.Queue)-1]
	}
}

//Set ...
func (lru *LRU) Set(key string, value ToDo) {
	if _, exist := lru.ToDoItems[key]; exist {
		index := getSliceIndex(key, lru.Queue...)
		if index > -1 {
			lru.setElementOnHead(key, index)
		}

		return
	}
	if lru.Capacity == len(lru.Queue) {
		lru.purge()
	}
	newToDo := ToDo{
		Name: value.Name,
		Done: value.Done,
	}
	lru.Lock()
	defer lru.Unlock()
	lru.ToDoItems[key] = newToDo
	lru.Queue = append([]string{key}, lru.Queue...)
}
