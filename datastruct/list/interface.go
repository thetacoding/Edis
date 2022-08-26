package list

type Expected func(a interface{}) bool

type Consumer func(i int, v interface{}) bool

type List interface {
	// Redis like API
	Create() *LinkedList	// todo 参考redis与godis决定是否实现该方法
	Empty()
	Release()
	AddNodeHead(val interface{})
	AddNodeTail(val interface{}) *LinkedList
	InsertNode(index int, val interface{})	// todo 参考redis实现
	DelNode(index int)
	GetIterator(direction int) *listIter
	Next() *node
	ReleaseIterator()
	Duplicate() *LinkedList
	SearchKey(val interface{}) *node	// todo
	Index(index int) *node
	Rewind() *listIter
	RewindTail() *listIter
	RewindTailToHead()
	RewindHeadToTail()
	Join(list *LinkedList)

	// Godis API
	Add(val interface{})
	Get(index int) (val interface{})
	Set(index int, val interface{})
	Insert(index int, val interface{})
	Remove(index int) (val interface{})
	RemoveLast() (val interface{})
	RemoveAllByVal(expected Expected) int
	RemoveByVal(expected Expected, count int) int
	ReverseRemoveByVal(expected Expected, count int) int
	Len() int
	ForEach(consumer Consumer)
	Contains(expected Expected) bool
	Range(start int, stop int) []interface{}
}
