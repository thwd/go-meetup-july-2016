package labels

import (
	"grammar"
)

type Label interface {
	Id() Id
	label() // private interface
}

var (
	// compile-time implementation checks
	_ Label = Byte{}
	_ Label = ByteSet{}
	_ Label = Any{}
	_ Label = Sequence{}
	_ Label = Alternation{}
	_ Label = Repetition{}
	_ Label = Option{}
	_ Label = &Recursion{}
)

type Id int

type Byte struct {
	id          Id
	Byte        byte
	Match, Fail Id
}

type ByteSet struct {
	id          Id
	ByteSet     grammar.ByteSet
	Match, Fail Id
}

type Any struct {
	id          Id
	Match, Fail Id
}

type Sequence struct {
	id            Id
	First, Second Id
}

type Alternation struct {
	id         Id
	This, That Id
}

type Repetition struct {
	id   Id
	Node Id
}

type Option struct {
	id   Id
	Node Id
}

type Recursion struct {
	id          Id
	Match, Fail Id
}

// label to go to when matching
// finished successfully
type Match struct {
	id Id
}

func NewMatch(i Id) Match {
	return Match{i}
}

// label to go to when matching
// failed
type Fail struct {
	id Id
}

func NewFail(i Id) Fail {
	return Fail{i}
}

// implement Label
func (Match) label()       {}
func (Fail) label()        {}
func (Byte) label()        {}
func (Sequence) label()    {}
func (Alternation) label() {}
func (Repetition) label()  {}
func (Any) label()         {}
func (Option) label()      {}
func (ByteSet) label()     {}
func (*Recursion) label()  {}

func (l Match) Id() Id       { return l.id }
func (l Fail) Id() Id        { return l.id }
func (l Byte) Id() Id        { return l.id }
func (l Sequence) Id() Id    { return l.id }
func (l Alternation) Id() Id { return l.id }
func (l Repetition) Id() Id  { return l.id }
func (l Option) Id() Id      { return l.id }
func (l Any) Id() Id         { return l.id }
func (l ByteSet) Id() Id     { return l.id }
func (l *Recursion) Id() Id  { return l.id }
