package chip8

type Stack struct {
	innerStack []int
}

func newStack() *Stack {
	stack := Stack{}
	stack.innerStack = []int{}
	return &stack
}

func (s *Stack) push(addr int) {
	s.innerStack = append(s.innerStack, addr)
}
