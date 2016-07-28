package grammar

func Any() Node {
	return AnyByte{}
}

func Alternate(ns ...Node) Node {
	if len(ns) == 0 {
		panic("")
	}
	ns = flattenAllAlternations(ns)
	ns = uniteSubsequentByteSets(ns)
	if len(ns) == 1 {
		return ns[0]
	}
	if IsNullable(ns[0]) {
		return ns[0]
	}
	return Alternation{ns[0], Alternate(ns[1:]...)}
}

func Sequenciate(ns ...Node) Node {
	if len(ns) == 0 {
		panic("")
	}
	if len(ns) == 1 {
		return ns[0]
	}
	return Sequence{ns[0], Sequenciate(ns[1:]...)}
}

func Repeat(n Node) Node {
	if _, k := n.(Repetition); k {
		return n
	}
	if t, k := n.(Option); k {
		return Repetition{t.Node}
	}
	return Repetition{n}
}

func Optionally(n Node) Node {
	if _, k := n.(Option); k {
		return n
	}
	if _, k := n.(Repetition); k {
		return n
	}
	return Option{n}
}

func Literal(s string) Node {
	ns := make([]Node, len(s), len(s))
	for i, l := 0, len(s); i < l; i++ {
		ns[i] = Byte(s[i])
	}
	return Sequenciate(ns...)
}

func Bytes(bs ...byte) ByteSet {
	s := ByteSet{}
	for _, b := range bs {
		s[b] = true
	}
	return s
}

func NotBytes(bs ...byte) ByteSet {
	s := AllBytes()
	for _, b := range bs {
		s[b] = false
	}
	return s
}

func ByteRange(from, to byte) ByteSet {
	if to < from {
		panic("")
	}
	s := ByteSet{}
	for b := from; b <= to; b++ {
		s[b] = true
	}
	return s
}

func NotByteRange(from, to byte) ByteSet {
	if to < from {
		panic("")
	}
	s := AllBytes()
	for b := from; b <= to; b++ {
		s[b] = false
	}
	return s
}

func IsNullable(n Node) bool {
	if _, k := n.(Option); k {
		return true
	}
	if _, k := n.(Repetition); k {
		return true
	}
	if t, k := n.(Sequence); k {
		return IsNullable(t.First) && IsNullable(t.Second)
	}
	if t, k := n.(Alternation); k {
		return IsNullable(t.That) // t.This will never be nullable
	}
	return false
}

func flattenAllAlternations(ns []Node) []Node {
	ms := make([]Node, 0, len(ns))
	for _, n := range ns {
		if _, k := n.(Alternation); k {
			ms = append(ms, flattenAlternation(n)...)
		} else {
			ms = append(ms, n)
		}
	}
	return ms
}

func flattenAlternation(n Node) []Node {
	if t, k := n.(Alternation); !k {
		return []Node{n}
	} else {
		return append(flattenAlternation(t.This), flattenAlternation(t.That)...)
	}
}

func uniteSubsequentByteSets(ns []Node) []Node {
	ms := make([]Node, 0, len(ns))
	for i, l := 0, len(ns); i < l; i++ {
		n, k := ns[i].(ByteSet)
		if !k || len(ms) == 0 {
			ms = append(ms, ns[i])
			continue
		}
		if m, k := ms[len(ms)-1].(ByteSet); k {
			ms[len(ms)-1] = m.Union(n)
		}
	}
	return ms
}
