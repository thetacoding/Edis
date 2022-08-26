package list

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
)

func ToString(list *LinkedList) string {
	arr :=make([]string, list.size)
	list.ForEach(func(i int, v interface{}) bool {
		integer, _ :=v.(int)
		arr[i] = strconv.Itoa(integer)
		return true
	})
	return "[" + strings.Join(arr, ", ") + "]"
}

func EchoList(list *LinkedList) {
	fmt.Println("list is : " + ToString(list) + ", list size is : " + strconv.Itoa(list.size))
}

func Catch() {
	if err := recover(); err != any(nil) {
		log.Println("恐慌异常：", err)
		return
	} else {
		panic(any("no expected panic"))
	}
}

func TestCreate(t *testing.T) {
	list := Create()
	fmt.Println("list head is : ", list.head)
	fmt.Println("list tail is : ", list.tail)
	fmt.Println("list size is : ", list.size)
}

func TestLinkedList_Empty(t *testing.T) {
	list := Create()
	for i := 0; i < 10; i++ {
		list.AddNodeTail(i)
	}
	fmt.Println("list is : ", ToString(list))
	list.Empty()
	fmt.Println("after executing the Empty() method : ", ToString(list))
}

func TestLinkedList_AddNodeHead(t *testing.T) {
	list := Create()
	for i := 0; i < 10; i++ {
		list.AddNodeTail(i)
	}
	fmt.Println("list is : ", ToString(list))
	list.ForEach(func(i int, v interface{}) bool {
		intVal, _ := v.(int)
		if intVal != i {
			t.Error("add test fail: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
		}
		return true
	})
}

func TestLinkedList_AddNodeTail(t *testing.T) {
	list := Create()
	for i := 0; i < 10; i++ {
		list.AddNodeTail(i)
	}
	fmt.Println("list is : ", ToString(list))
	list.ForEach(func(i int, v interface{}) bool {
		intVal, _ := v.(int)
		if intVal != i {
			t.Error("add test fail: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(intVal))
		}
		return true
	})
}

func TestLinkedList_InsertNode(t *testing.T) {
	list := Create()
	// 索向前越界
	f1 := func() {
		defer Catch()
		list.InsertNode(-1, 1)
	}
	f1()
	// 索引向后越界
	f2 := func() {
		defer Catch()
		list.InsertNode(2, 2)
	}
	f2()
	// 索引刚好指向列尾
	list.InsertNode(0, 0)
	list.InsertNode(1, 1)
	list.InsertNode(2, 2)
	fmt.Println(ToString(list))
	// 插入到列表头
	list.InsertNode(0, 10)
	fmt.Println(ToString(list))
	// 插入列表中间
	list.InsertNode(1, 11)
	fmt.Println(ToString(list))
}

func TestLinkedList_DelNode(t *testing.T) {
	list := Create()
	for i := 0; i < 10; i++ {
		list.AddNodeTail(i)
	}
	EchoList(list)
	f1 := func() {
		defer Catch()
		list.DelNode(-1)
	}
	f1()
	f2 := func() {
		defer Catch()
		list.DelNode(10)
	}
	f2()
	list.DelNode(4)
	EchoList(list)
	list.DelNode(0)
	EchoList(list)
	list.DelNode(list.size-1)
	EchoList(list)
}































