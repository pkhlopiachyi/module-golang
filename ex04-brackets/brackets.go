package brackets

type Stack struct {
	size int
	top  *elem
}

type elem struct {
	val  rune
	next *elem
}

func New() *Stack {
	return &Stack{0, nil}

}

func (stack *Stack) Push(inter rune) {

	value := &elem{inter, stack.top}
	stack.top = value
	stack.size++
}

func (stack *Stack) Pop() rune {
	if stack.size > 0 {
		value := stack.top
		stack.top = value.next
		stack.size--
		return value.val
	}
	return 0
}

func (stack *Stack) Len() int {
	return stack.size
}

func Bracket(bracketsLine string) (bool, error) {
	var stack *Stack = New()
	var closeBr rune
	for _, bracket := range bracketsLine {
		switch {
		case bracket == 123:
			stack.Push(125)
		case bracket == 40:
			stack.Push(41)
		case bracket == 91:
			stack.Push(93)
		default:
			closeBr = stack.Pop()
			if closeBr != bracket {
				return false, nil
			}
		}
	}
	return stack.Len() == 0, nil
}
