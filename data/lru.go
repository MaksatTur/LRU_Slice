package data

import (
	"fmt"
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

func (lru *LRU) setElementOnHead(key string) {
	lru.Queue = append([]string{key}, lru.Queue...)
	fmt.Println(lru.Queue)
	// slice := make([]string, 0)
	// slice = append(slice, key)
	// fmt.Println(slice)
	// slice = append(slice, lru.Queue...)
	// fmt.Println(slice)
	// copy(*lru.Queue, &slice)

	// return slice
	// tmp := lru.Queue[index]
	// lru.Queue[index] = lru.Queue[0]
	// lru.Queue[0] = tmp
}

func (lru *LRU) purge() {
	if len(lru.Queue) > 0 {
		fmt.Println(lru.Queue[len(lru.Queue)-1])
		delete(lru.ToDoItems, lru.Queue[len(lru.Queue)-1])
		lru.Queue = lru.Queue[:len(lru.Queue)-1]
	}
}

//Set ...
func (lru *LRU) Set(key string, value ToDo) {
	if _, exist := lru.ToDoItems[key]; exist {
		// lock
		lru.Lock()
		defer lru.Unlock()
		//index := getSliceIndex(key, lru.Queue...)
		//lru.setElementOnHead(key, index)
		//copy(lru.Queue, lru.setElementOnHead(key))
		lru.setElementOnHead(key)
		return
	}
	if lru.Capacity == len(lru.Queue) {
		//fmt.Println("full")
		lru.purge()
		//fmt.Println("done")
	}
	newToDo := ToDo{
		Name: value.Name,
		Done: value.Done,
	}
	// fmt.Printf("Key = %s\n", key)
	lru.Lock()
	defer lru.Unlock()
	lru.ToDoItems[key] = newToDo
	// lru.Queue = append(lru.Queue, key) /// тут
	//slice := lru.setElementOnHead(key)
	//fmt.Printf("slice = %+v\n", slice)
	//copy(lru.Queue, slice)
	//lru.setElementOnHead(key)
	// set element on head
	lru.Queue = append([]string{key}, lru.Queue...)
	fmt.Printf("Queue = %+v\n", lru.Queue)
	//index := getSliceIndex(key, lru.Queue...)
	// if index > -1 {
	// 	lru.setElementOnHead(key, index)
	// }
	//fmt.Printf("Queue = %+v\n", lru.Queue)
	// fmt.Printf("ToDos = %+v\n", lru.ToDoItems)
}
