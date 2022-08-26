package list

import "Edis/datastruct/errs"

// LinkedList 双链接列表
type LinkedList struct {
	head *node
	tail *node
	size int
}

// node 列表节点
type node struct {
	val interface{}
	prev *node
	next *node
}

// list 迭代器
type listIter struct {
	next *node
	reverse bool
}

/*
	LinkedList 实现 List 接口方法
 */

// Create 初始化一个空链表并返回
func Create() *LinkedList {
	list := &LinkedList{
		head:nil,
		tail:nil,
		size:0,
	}
	return list
}

// Empty 清空列表所有节点
func (list *LinkedList) Empty() {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	current := list.head
	size := list.size
	var next *node
	for size > 0 {
		next = current.next
		current.next = nil	// for gc
		if next != nil {
			next.prev = nil		// for gc
		}
		current = next
		size--
	}
	list.head = nil
	list.tail = nil
	list.size = 0
}

// Release 释放列表 todo go中没必要实现
// 不能使用"func (list *LinkedList) Release() {...}"的原因是这里的list只是方法内一个指向LinkedList的指针变量
// list被赋值nil并不影响调用Release方法的那个指向LinkedList的指针变量或值变量
func Release(list *LinkedList) {
	panic(any("not supported"))
}

// AddNodeHead 列表添加元素到列头
func (list *LinkedList) AddNodeHead(val interface{})  {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	n := &node{
		val:  val,
	}
	if list.head == nil {
		// 空列表
		list.head = n
		list.tail = n
	} else {
		n.next = list.head
		list.head.prev = n
		list.head = n
	}
	list.size++
}

// AddNodeTail 列表添加元素到列尾
func (list *LinkedList) AddNodeTail(val interface{})  {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	n := &node{
		val:  val,
	}
	if list.tail == nil {
		// 空列表
		list.head = n
		list.tail = n
	} else {
		n.prev = list.tail
		list.tail.next = n
		list.tail = n
	}
	list.size++
}

// InsertNode 列表指定位置插入元素
func (list *LinkedList) InsertNode(index int, val interface{}) {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	if index < 0 || index > list.size {
		panic(any(errs.INDEX_OUT_BOUND))
	}
	// 刚好要插入队尾
	if index == list.size {
		list.Add(val)
		return
	}
	pivot := list.find(index)
	n := &node{
		val:  val,
		prev: pivot.prev,
		next: pivot,
	}
	if pivot.prev == nil {
		list.head = n
	} else {
		pivot.prev.next = n
	}
	pivot.prev = n
	list.size++
}

// DelNode 列表删除指定位置元素
func (list *LinkedList) DelNode(index int) {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	if index < 0 || index >= list.size {
		panic(any(errs.INDEX_OUT_BOUND))
	}
	pivot := list.find(index)
	if pivot.prev == nil {
		list.head = pivot.next
	} else {
		pivot.prev.next = pivot.next
	}
	if pivot.next == nil {
		list.tail = pivot.prev
	} else {
		pivot.next.prev = pivot.prev
	}
	list.size--
	// for gc
	pivot.prev = nil
	pivot.next = nil
}

// GetIterator 获取列表 list 的迭代器，每次调用 listNext() 方法返回下一个元素
func (list *LinkedList) GetIterator(reverse bool) *listIter {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	iter := &listIter{}
	if reverse {
		iter.next = list.tail
	} else {
		iter.next = list.head
	}
	iter.reverse = reverse
	return iter
}

// 释放 list 迭代器
func (iter *listIter) ReleaseIterator() {
	iter.next = nil
}

// 重置迭代器到列表头
func (iter *listIter) listRewind(list *LinkedList) {
	iter.next = list.head
	iter.reverse = false
}

// 重置迭代器到列表尾
func (iter *listIter) listRewindTail(list *LinkedList) {
	iter.next = list.tail
	iter.reverse = true
}

// 返回列表下一个元素
func (iter *listIter) listNext() *node {
	current := iter.next
	if current != nil {
		if iter.reverse {
			iter.next = current.prev
		} else {
			iter.next = current.next
		}
	}
	return current
}

// 复制整个列表
func (list *LinkedList) listDup() *LinkedList {

}















// Add 列表添加元素
func (list *LinkedList) Add(val interface{})  {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	n := &node{
		val:  val,
	}
	if list.tail == nil {
		// 空列表
		list.head = n
		list.tail = n
	} else {
		n.prev = list.tail
		list.tail.next = n
		list.tail = n
	}
	list.size++
}

// find 获取列表中指定索引的元素
// 该方法假定 list 不为 nil，index 符合索引范围
func (list *LinkedList) find(index int) (n *node) {
	if index < list.size/2 {
		n = list.head
		for i := 0; i < index; i++ {
			n = n.next
		}
	} else {
		n = list.tail
		for i := list.size - 1; i > index; i-- {
			n = n.prev
		}
	}
	return n
}

// get 获取列表中指定索引的元素的值
func (list *LinkedList) Get(index int) (val interface{}) {
	return list.find(index).val
}

// Insert 列表指定索引位置插入元素
func (list *LinkedList) Insert(index int, val interface{})  {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	if index < 0 || index > list.size {
		panic(any(errs.INDEX_OUT_BOUND))
	}
	// 刚好要插入队尾
	if index == list.size {
		list.Add(val)
		return
	}
	pivot := list.find(index)
	// 如果 list 为空，且 index 不为 0，则 list.find(index) 会报错
	// 程序走到此处，可以表面 list 不为空
	n := &node{
		val:  val,
		prev: pivot.prev,
		next: pivot,
	}
	if pivot.prev == nil {
		list.head = n
	} else {
		pivot.prev.next = n
	}
	pivot.prev = n
	list.size++
}

// removeNode 移除列表中节点
func (list *LinkedList) removeNode(n *node) {
	if n.prev == nil {
		list.head = n.next
	} else {
		n.prev.next = n.next
	}
	if n.next == nil {
		list.tail = n.prev
	} else {
		n.next.prev = n.prev
	}

	// for gc
	n.prev = nil
	n.next = nil

	list.size--
}

// Remove 移除列表中指定索引的节点
func (list *LinkedList) Remove(index int) (val interface{}) {
	n := list.find(index)
	list.removeNode(n)
	return n.val
}

// Removetail 移除列表最后一个节点
func (list *LinkedList) Removetail() (val interface{}) {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	if list.tail == nil {
		// empty list
		return nil
	}
	n := list.tail
	list.removeNode(n)
	return n.val
}

func (list *LinkedList) RemoveAllByVal(expected Expected) int {
	if list == nil {
		panic(any(errs.INDEX_OUT_BOUND))
	}
	n := list.head
	removed := 0
	var nextNode *node
	for n != nil {
		nextNode = n.next
		if expected(n.val) {
			list.removeNode(n)
			removed++
		}
		n = nextNode
	}
	return removed
}

func (list *LinkedList) ForEach(consumer Consumer) {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	n := list.head
	i := 0
	for n != nil {
		goNext := consumer(i, n.val)
		if !goNext {
			break
		}
		i++
		n = n.next
	}
}




























