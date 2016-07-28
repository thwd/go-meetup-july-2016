package grammar

type Node interface {
	node() // private interface
}

var (
	// compile-time implementation checks
	_ Node = Byte(0)
	_ Node = AnyByte{}
	_ Node = Sequence{}
	_ Node = Alternation{}
	_ Node = Repetition{}
	_ Node = Option{}
	_ Node = &Recursion{}
	_ Node = ByteSet{}
)

// matches a specific byte
type Byte byte

// consumes and matches any one byte
type AnyByte struct{}

// matches against First, then against Second
type Sequence struct {
	First, Second Node
}

// tries to match This and if it fails, tries That
// (behaves like PEG alternation)
type Alternation struct {
	This, That Node
}

// matches Node zero or more times greedily
// (will never not match)
type Repetition struct {
	Node Node
}

// matches Node zero or one time
// (will never not match)
type Option struct {
	Node Node
}

// used to construct recursive grammar trees
type Recursion struct {
	Node Node
}

func Recurse(f func(*Recursion) Node) *Recursion {
	r := &Recursion{nil}
	r.Node = f(r)
	return r
}

// matches any one byte that has a decimal value
// equivalent to any of the indices that are true
type ByteSet [256]bool

// constructs a ByteSet equivalent to AnyByte
// useful for building not-sets
func AllBytes() ByteSet {
	s := ByteSet{}
	for i, _ := range s {
		s[i] = true
	}
	return s
}

// self explanatory
func (s ByteSet) Union(t ByteSet) ByteSet {
	for i, _ := range t {
		s[i] = t[i] || s[i]
	}
	return s
}

// self explanatory
func (s ByteSet) Intersection(t ByteSet) ByteSet {
	for i, _ := range t {
		s[i] = t[i] && s[i]
	}
	return s
}

// returns the inclusive range boundaries
// in which a ByteSet is true
func (s ByteSet) Ranges() [][2]byte {
	rs := make([][2]byte, 0, 8)
	j := -1
	for i, _ := range s {
		if s[i] && j == -1 {
			j = i
		} else if !s[i] && j != -1 {
			rs = append(rs, [2]byte{byte(j), byte(i - 1)})
			j = -1
		}
	}
	if j != -1 {
		rs = append(rs, [2]byte{byte(j), 255})
	}
	return rs
}

// returns a slice of all the bytes that
// are set in the ByteSet
func (s ByteSet) Bytes() []byte {
	bs := make([]byte, 0, len(s))
	for i, _ := range s {
		if s[i] {
			bs = append(bs, byte(i))
		}
	}
	return bs
}

// private implementations
func (Byte) node()        {}
func (Sequence) node()    {}
func (Alternation) node() {}
func (Repetition) node()  {}
func (Option) node()      {}
func (AnyByte) node()     {}
func (ByteSet) node()     {}
func (*Recursion) node()  {}
