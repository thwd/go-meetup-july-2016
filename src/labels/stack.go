package labels

type LabelStack struct {
	Labels []Label
}

func (q *LabelStack) Push(n Label) {
	q.Labels = append(q.Labels, n)
}

// might panic
func (q *LabelStack) Pop() Label {
	n := q.Peek()
	q.Labels = q.Labels[:len(q.Labels)-1]
	return n
}

// might panic
func (q *LabelStack) Peek() Label {
	return q.Labels[len(q.Labels)-1]
}
