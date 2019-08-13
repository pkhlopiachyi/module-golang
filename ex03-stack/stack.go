package stack

type Stack struct {
	size int
	top  *elem
}

type elem struct {
	val  int
	next *elem
}

func New() *Stack {
	return &Stack{}
}

func (stack *Stack) Push(value int) {
	stack.top = &elem{value, stack.top}
	stack.size++
}

func (stack *Stack) Pop() int {
	var value int
	if stack.size > 0 {
		value = stack.top.val
		stack.top = stack.top.next
		stack.size--
		return value
	}
	return -1
}
