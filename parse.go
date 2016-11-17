package boolean

import (
	"fmt"
	"log"
	"strings"
)

func Precedence(item item) int {
	switch item.typ {
	case itemOr:
		return 1
	case itemAnd:
		return 2
	default:
		panic(fmt.Sprintf("unknown item type: %s", item.typ))
	}
	return 0
}

type Node interface {
}

type Identifier struct {
	ident string
}

type UnaryOperator struct {
	item item
	rhs  Node
}

type Operator struct {
	item     item
	lhs, rhs Node
}

type ParenNode struct {
	Node Node
}

type Tree struct {
	lex    *lexer
	text   string
	Root   Node
	tokens []item
}

func New() *Tree {
	return &Tree{}
}

func (t *Tree) Parse(text string) (tree *Tree, err error) {
	defer t.recover(&err)
	t.lex = lex(text)
	t.Root = t.parseExpr()
	return t, nil
}

func (t *Tree) recover(err *error) {
	e := recover()
	if e != nil {
		t.lex.drain()
		*err = e.(error)
	}
	return
}

func (t *Tree) next() item {
	if len(t.tokens) == 0 {
		return t.lex.nextItem()
	}
	token := t.tokens[len(t.tokens)-1]
	t.tokens = t.tokens[:len(t.tokens)-1]
	return token
}

func (t *Tree) backup(token item) {
	t.tokens = append(t.tokens, token)
}

func (t *Tree) parseExpr() Node {
	root := &Operator{rhs: t.parseUnaryExpr()}
	// now we start parsing pairs of binary operator and its rhs
	for {
		switch item := t.next(); item.typ {
		case itemAnd, itemOr:
			op := &Operator{item: item}
			// 1. find the RHS of the binary op which must be an expr
			rhs := t.parseUnaryExpr()
			// 2. add this binary operator to the tree, see influxsql parser algorithm
			for node := root; ; {
				r, ok := node.rhs.(*Operator)
				if !ok || Precedence(r.item) >= Precedence(op.item) {
					op.rhs = rhs
					op.lhs = node.rhs
					node.rhs = op
					break
				}
				node = r
			}
		case itemRightParen:
			t.backup(item)
			return root.rhs
		case itemEOF:
			return root.rhs
		case itemError:
			panic(fmt.Errorf("%s", item))
		default:
			panic(fmt.Errorf("unexpected %s", item))
		}
	}
	return nil
}

func (t *Tree) parseUnaryExpr() Node {
	item := t.next()
	switch item.typ {
	case itemIdentifier:
		// adauga-l ca LHS
		return &Identifier{item.val}
	case itemLeftParen:
		// parse an expression that must end with right paren
		node := t.parseExpr()
		// make sure that we end with right paren
		if item := t.next(); item.typ != itemRightParen {
			panic(fmt.Errorf("unclosed right paren: unexpected %s", item))
		}
		return &ParenNode{Node: node}
	case itemNot:
		return &UnaryOperator{item: item, rhs: t.parseUnaryExpr()}
	case itemEOF:
		panic(fmt.Errorf("unexpected EOF"))
	case itemError:
		panic(fmt.Errorf("%s", item))
	default:
		panic(fmt.Errorf("unexpected %s", item))
	}
}

// PrintNode prints a node
// only for debugging
func PrintNode(node Node, count int) string {
	switch el := node.(type) {
	case *Identifier:
		return fmt.Sprintf("Name: %s", el.ident)
	case *UnaryOperator:
		return fmt.Sprintf("Name: %s, Children:\n%s%s", el.item.val, strings.Repeat("\t", count), PrintNode(el.rhs, count+1))
	case *Operator:
		return fmt.Sprintf("Name: %s, Children:\n%s%s\n%s%s",
			el.item.val, strings.Repeat("\t", count), PrintNode(el.lhs, count+1), strings.Repeat("\t", count), PrintNode(el.rhs, count+1))
	case *ParenNode:
		return fmt.Sprintf("(%s)", PrintNode(el.Node, count))
	default:
		log.Fatalf("unknown node type: %T", el)
	}
	return ""
}
