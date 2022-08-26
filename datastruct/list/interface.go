package list

type Expected func(a interface{}) bool

type Consumer func(i int, v interface{}) bool

type List interface {
	// Redis like API
	Create() *LinkedList
	Empty()
	Release()
	AddNodeHead(val interface{})
	AddNodeTail(val interface{}) *LinkedList
	InsertNode(index int, val interface{})
	DelNode(index int)
	GetIterator(direction int) *listIter
	Next() *listNode
	ReleaseIterator()
	Duplicate() *LinkedList
	SearchKey(val interface{}) *listNode	// todo
	Index(index int) *listNode
	Rewind() *listIter
	RewindTail() *listIter
	RewindTailToHead()
	RewindHeadToTail()
	Join(list *LinkedList)

	// Godis like API
	ForEach(consumer Consumer)
}
