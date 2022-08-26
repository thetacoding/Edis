package list

import (
	"Edis/datastruct/errs"
)

// LinkedList 双链接列表
type LinkedList struct {
	head *listNode
	tail *listNode
	size int
	dup func(node *listNode) *interface{}
	free func(node *listNode)
	match func(node *listNode) bool
}

// listNode 列表节点
type listNode struct {
	val interface{}
	prev *listNode
	next *listNode
}

// listIter 迭代器
type listIter struct {
	next *listNode
	reverse bool
}


/* Create a new list. Return the pointer to the new list.
 * 创建列表，返回列表指针 */
func Create() *LinkedList {
	list := &LinkedList{
		head:nil,
		tail:nil,
		size:0,
		dup:nil,
		match:nil,
	}
	return list
}

/* Remove all the elements from the list without destroying the list itself.
 * 清空列表所有元素 */
func (list *LinkedList) Empty() {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	current := list.head
	size := list.size
	var next *listNode
	for size > 0 {
		next = current.next
		current.next = nil	// for gc
		if next != nil {
			next.prev = nil		// for gc
		}
		if list.free != nil {
			list.free(current)
		}
		current = next
		size--
	}
	list.head = nil
	list.tail = nil
	list.size = 0
}

/* Free the whole list.
 * This function is equivalent to Empty(), only for consistency with redis API.
 * 释放列表，等同于Empty()方法，仅为与redis接口保持一致 */
func (list *LinkedList) Release() {
	list.Empty()
}

/* Add a new listNode to the list, to head, containing the specified 'value' pointer as value.
 * 添加元素到列表头部 */
func (list *LinkedList) AddNodeHead(val interface{})  {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	n := &listNode{
		val:  val,
	}
	if list.head == nil {
		// empty list
		list.head = n
		list.tail = n
	} else {
		n.next = list.head
		list.head.prev = n
		list.head = n
	}
	list.size++
}

/* Add a new listNode to the list, to tail, containing the specified 'value' pointer as value.
 * 添加元素到列表尾部 */
func (list *LinkedList) AddNodeTail(val interface{})  {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	n := &listNode{
		val:  val,
	}
	if list.tail == nil {
		// empty list
		list.head = n
		list.tail = n
	} else {
		n.prev = list.tail
		list.tail.next = n
		list.tail = n
	}
	list.size++
}

/* Insert a new element to the specified zero-based index.
 * 插入元素到指定位置 */
func (list *LinkedList) InsertNode(index int, val interface{}) {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	if index < 0 || index > list.size {
		panic(any(errs.INDEX_OUT_BOUND))
	}
	// right at the tail
	if index == list.size {
		list.AddNodeTail(val)
		return
	}
	pivot := list.Index(index)
	n := &listNode{
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

/* Remove the specified listNode from the specified list.
 * It's up to the caller to free the private value of the node.
 * 删除指定位置元素 */
func (list *LinkedList) DelNode(index int) {
	if list == nil {
		panic(any(errs.NIL_LIST))
	}
	if index < 0 || index >= list.size {
		panic(any(errs.INDEX_OUT_BOUND))
	}
	pivot := list.Index(index)
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

/* Returns a list iterator 'iter'. After the initialization every
 * call to listNext() will return the next element of the list.
 * 返回一个迭代器，初始化后每次调用listNext()会返回列表的下一个元素 */
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

/* Release the iterator memory.
 * Discard iter through setting the private next of the listIter. Only for consistency with redis API
 * 释放迭代器。
 * 令迭代器的next指向nil来废弃迭代器。这只是为了与redis的API保持一致 */
func (iter *listIter) ReleaseIterator() {
	iter.next = nil
}

/* Rewind the iterator to the head.
 * 迭代器倒回列表头部 */
func (iter *listIter) listRewind(list *LinkedList) {
	iter.next = list.head
	iter.reverse = false
}

/* Rewind the iterator to the tail.
 * 迭代器倒回列表尾部 */
func (iter *listIter) listRewindTail(list *LinkedList) {
	iter.next = list.tail
	iter.reverse = true
}

/* Return the next element of the list.
 * 返回列表下一个元素 */
func (iter *listIter) listNext() *listNode {
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

/* Copy the whole list.
 * The 'dup' method will be performed to copy the listNode if set.
 * 复制整个列表。如果有设置，'dup'函数将会被使用 */
func (list *LinkedList) Duplicate() *LinkedList {
	copy :=  Create()
	copy.dup = list.dup
	copy.free = list.free
	copy.match = list.match
	iter := list.GetIterator(false)

	node := iter.listNext()
	for node != nil {
		var val *interface{}
		if copy.dup != nil {
			val = copy.dup(node)
			if val == nil {
				copy.Release()
				return nil
			}
		} else {
			val = &node.val
		}
		// redis源码实现中这里会判断列表添加元素是否成功，因为其源码实现中考虑了zmalloc分配内存空间失败的情形
		copy.AddNodeTail(*val)
		node = iter.listNext()
	}
	return copy
}

/* Search the list for a node matching a given key
 * The 'match' method will be performed to compare the listNode if set.
 * On success the first matching node pointer is returned. If no matching node exists nil is returned.
 * 如果设置了list的match函数，则使用match函数进行匹配，否则比较节点的val值。
 * 返回第一个匹配成功的节点，如果没有可匹配的节点则返回nil。
 */
func (list *LinkedList) SearchKey(key interface{}) *listNode {
	iter := list.GetIterator(false)
	node := iter.listNext()
	for node != nil {
		if list.match != nil {
			if list.match(node) {
				return node
			}
		} else {
			if node.val == key {
				return node
			}
		}
		node = iter.listNext()
	}
	return nil
}

/* Return the element at the specified zero-based index
 * where 0 is the head, 1 is the element next to head
 * and so on.Negative integers are used in order to count
 * from the tail, -1 is the last element, -2 the penultimate
 * and so on. If the index is out of range nil is returned.
 * 返回指定索引的元素，索引从0开始。索引为负数表示从列尾开始向前
 * 搜索，-1表示列尾元素，-2表示列尾元素的前一个元素，...。索引超出
 * 界限将返回nil。 */
func (list *LinkedList) Index(index int) *listNode {
	var node *listNode;
	if index < 0 {
		index = -index-1
		node = list.tail
		for index > 0 && node != nil {
			node = node.prev
			index--
		}
	} else {
		node = list.head
		for index > 0 && node != nil {
			node = node.next
			index--
		}
	}
	return node
}

/* The 'Consumer' method will be performed on each listNode.
 * 每一个节点都将执行'Consumer'方法 */
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