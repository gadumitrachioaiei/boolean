package boolean

import "fmt"

func (t *Tree) Execute(data map[string]bool) (val bool, err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			err = err1.(error)
		}
	}()
	return execute(t.Root, data), nil
}

func execute(node Node, data map[string]bool) bool {
	switch op := node.(type) {
	case *Identifier:
		if value, ok := data[op.ident]; ok {
			return value
		}
		panic(fmt.Errorf("undefined identifier: %s", op.ident))
	case *Operator:
		switch op.item.typ {
		case itemAnd:
			return execute(op.lhs, data) && execute(op.rhs, data)
		case itemOr:
			return execute(op.lhs, data) || execute(op.rhs, data)
		default:
			panic(fmt.Errorf("unknown unary operator: %s", op.item))
		}
	case *UnaryOperator:
		if op.item.typ == itemNot {
			return !execute(op.rhs, data)
		}
		panic(fmt.Errorf("unknown unary operator: %s", op.item))
	case *ParenNode:
		return execute(op.Node, data)
	default:
		panic(fmt.Errorf("unknown node type: %v", op))
	}
}
