package labels

import (
	"fmt"
	"grammar"
)

func Labelize(i func() Id, n grammar.Node) []Label {
	p := make([]Label, 0, 128)
	q := LabelStack{}
	traverseNode(n, func(x grammar.Node) {
		l := (Label)(nil)
		switch t := x.(type) {
		case grammar.ByteSet:
			l = ByteSet{i(), t, -1, -1}
		case grammar.Byte:
			l = Byte{i(), (byte)(t), -1, -1}
		case grammar.AnyByte:
			l = Any{i(), -1, -1}
		case *grammar.Recursion:
			l = &Recursion{i(), -1, -1}
		case grammar.Sequence:
			a, b := q.Pop().Id(), q.Pop().Id()
			l = Sequence{i(), b, a}
		case grammar.Alternation:
			a, b := q.Pop().Id(), q.Pop().Id()
			l = Alternation{i(), b, a}
		case grammar.Repetition:
			l = Repetition{i(), q.Pop().Id()}
		case grammar.Option:
			l = Option{i(), q.Pop().Id()}
		default:
			panic(fmt.Sprintf("unhandled type %T\n", x))
		}
		q.Push(l)
		p = append(p, l)
	})
	return p
}

// traverse nodes in post-order
func traverseNode(n grammar.Node, f func(grammar.Node)) {
	switch t := n.(type) {
	case grammar.ByteSet:
	case grammar.AnyByte:
	case grammar.Byte:
	case *grammar.Recursion:
	case grammar.Sequence:
		traverseNode(t.First, f)
		traverseNode(t.Second, f)
	case grammar.Alternation:
		traverseNode(t.This, f)
		traverseNode(t.That, f)
	case grammar.Repetition:
		traverseNode(t.Node, f)
	case grammar.Option:
		traverseNode(t.Node, f)
	default:
		panic(fmt.Sprintf("unhandled type %T\n", n))
	}
	f(n)
}
