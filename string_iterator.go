package mob

type StringIterator interface {
	Len() int
	At(idx int) string
}

type StringArray []string

func (this StringArray) Len() int {
	return len(this)
}

func (this StringArray) At(idx int) string {
	return this[idx]
}
