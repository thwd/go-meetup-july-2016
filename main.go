package main

import (
	"fmt"
	"grammar"
	"labels"
	"os"
	"reflect"
	"text/template"
)

func main() {
	json := grammar.JsonGrammar()

	labelCounter := labels.Id(0)
	newLabel := func() labels.Id { labelCounter++; return labelCounter }

	labelStack := labels.Labelize(newLabel, json.(*grammar.Recursion).Node)

	match, fail := newLabel(), newLabel()
	entrypoint := LinkLabels(labelStack, match, fail)

	labelStack = append(labelStack, labels.NewMatch(match))
	labelStack = append(labelStack, labels.NewFail(fail))

	tt := template.New("parser.tmpl")

	tt.Funcs(template.FuncMap{
		"LabelType": func(l labels.Label) string {
			t := reflect.TypeOf(l)
			for t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			return t.Name()
		},
		"StructField": func(s string, i interface{}) interface{} {
			v := reflect.ValueOf(i)
			for v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			return v.FieldByName(s).Interface()
		},
		"ByteSetRanges": func(l labels.ByteSet) [][2]byte {
			return l.ByteSet.Ranges()
		},
		"ByteSetEmpty": func(l labels.ByteSet) bool {
			return len(l.ByteSet.Ranges()) == 0
		},
	})

	_, e := tt.ParseFiles("parser.tmpl")

	if e != nil {
		panic(e)
	}

	e = tt.Execute(os.Stdout, map[string]interface{}{
		"Labels":     labelStack,
		"EntryPoint": entrypoint,
	})

	if e != nil {
		panic(e)
	}
}

// m = match label
// f = fail label
func LinkLabels(ls []labels.Label, m, f labels.Id) labels.Id {
	x := len(ls) - 1
	switch t := ls[x].(type) {
	case labels.ByteSet:
		t.Match = m
		t.Fail = f
		ls[x] = t
		return t.Id()
	case labels.Byte:
		t.Match = m
		t.Fail = f
		ls[x] = t
		return t.Id()
	case labels.Any:
		t.Match = m
		t.Fail = f
		ls[x] = t
		return t.Id()
	case *labels.Recursion:
		t.Match = m
		t.Fail = f
		return t.Id()
	case labels.Sequence:
		s := LinkLabels(labelsUpTo(ls, t.Second), m, f)
		r := LinkLabels(labelsUpTo(ls, t.First), s, f)
		return r

	case labels.Alternation:
		s := LinkLabels(labelsUpTo(ls, t.That), m, f)
		r := LinkLabels(labelsUpTo(ls, t.This), m, s)
		return r
	case labels.Repetition:
		s := LinkLabels(labelsUpTo(ls, t.Node), t.Node, m)
		s = LinkLabels(labelsUpTo(ls, t.Node), s, m)
		return s

	case labels.Option:
		return LinkLabels(labelsUpTo(ls, t.Node), m, m)
	default:
		panic(fmt.Sprintf("%T\n", ls[x]))
	}
}

func labelsUpTo(ls []labels.Label, id labels.Id) []labels.Label {
	for i := len(ls) - 1; i >= 0; i-- {
		if ls[i].Id() == id {
			return ls[:i+1]
		}
	}
	panic("never reached")
}
