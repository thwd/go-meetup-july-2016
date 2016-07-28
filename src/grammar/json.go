package grammar

func JsonGrammar() Node {
	// -?(0|[1-9][0-9]*)(.[0-9]*(e|E)[+-]?[1-9][0-9*])?
	n := Sequenciate(
		Optionally(Byte('-')),
		Alternate(
			Byte('0'),
			Sequenciate(
				ByteRange('1', '9'),
				Repeat(ByteRange('0', '9')),
			),
		),
		Optionally(Sequenciate(
			Byte('.'),
			Repeat(ByteRange('0', '9')),
			Optionally(Sequenciate(
				Bytes('E', 'e'),
				Optionally(Bytes('+', '-')),
				ByteRange('1', '9'),
				Repeat(ByteRange('0', '9')),
			)),
		)),
	)
	h := Alternate(
		ByteRange('A', 'F'),
		ByteRange('a', 'f'),
		ByteRange('0', '9'),
	)

	_ = h

	c := ByteRange(' ', '~').Intersection(NotBytes('\\', '"'))

	s := Sequenciate(
		Byte('"'),
		Repeat(Alternate(
			Sequenciate(Byte('\\'), Alternate(
				Sequenciate(Byte('u'), h, h, h, h),
				Bytes('"', '\\', '/', 'b', 'f', 'n', 'r', 't'),
			)),
			Sequenciate(c, Repeat(c)),
		)),
		Byte('"'),
	)
	b := Alternate(
		Literal("true"),
		Literal("false"),
	)
	u := Literal("null")

	x := &Recursion{}

	a := Sequenciate(
		Byte('['),
		Optionally(Sequenciate(
			x,
			Repeat(Sequenciate(
				Byte(','),
				x,
			)),
		)),
		Byte(']'),
	)

	o := Sequenciate(
		Byte('{'),
		Optionally(Sequenciate(
			s, Byte(':'), x,
			Repeat(Sequenciate(
				Byte(','),
				s, Byte(':'), x,
			)),
		)),
		Byte('}'),
	)

	x.Node = Alternate(o, a, s, b, n, u)

	return x
}
