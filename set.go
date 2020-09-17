package mob

import "encoding/json"

type (
	Set interface {
		Len() int
		Add(v interface{}) bool
		Contains(v interface{}) bool
		Remove(v interface{})
		Collection() []interface{}
		Foreach(cb func(interface{}))
		Intersection(other Set) Set
		Subtraction(other Set) Set
		Clone() Set
	}

	setImpl struct {
		holder map[interface{}]uint8
	}
)

func NewSet() Set {
	return NewSetWithCap(16)
}

func NewSetWithCap(capHint int) Set {
	return &setImpl{
		holder: make(map[interface{}]uint8, capHint),
	}
}

func NewSetFromStrings(arr []string) Set {
	out := NewSetWithCap(len(arr))
	for _, it := range arr {
		out.Add(it)
	}
	return out
}

func (this *setImpl) Len() int {
	return len(this.holder)
}

func (this *setImpl) Add(v interface{}) (ok bool) {
	if this.Contains(v) {
		return false
	}
	this.holder[v] = 0
	return true
}

func (this *setImpl) Contains(v interface{}) bool {
	_, ok := this.holder[v]
	return ok
}

func (this *setImpl) Remove(v interface{}) {
	delete(this.holder, v)
}

func (this *setImpl) Collection() []interface{} {
	out := make([]interface{}, 0, len(this.holder))
	this.Foreach(func(it interface{}) {
		out = append(out, it)
	})
	return out
}

func (this *setImpl) Foreach(cb func(interface{})) {
	if cb == nil {
		return
	}
	for k := range this.holder {
		cb(k)
	}
}

func (this *setImpl) Intersection(other Set) Set {
	var smaller Set = this
	if this.Len() > other.Len() {
		smaller, other = other, this
	}
	outSet := NewSetWithCap(smaller.Len())
	smaller.Foreach(func(i interface{}) {
		if other.Contains(i) {
			outSet.Add(i)
		}
	})
	return outSet
}

func (this *setImpl) Subtraction(other Set) Set {
	in := this.Intersection(other)
	out := NewSet()
	this.Foreach(func(i interface{}) {
		if !in.Contains(i) {
			out.Add(i)
		}
	})
	other.Foreach(func(i interface{}) {
		if !in.Contains(i) {
			out.Add(i)
		}
	})
	return out
}

func (this *setImpl) Clone() Set {
	thisByte, _ := json.Marshal(this.holder)
	outSet := NewSet().(*setImpl)
	_ = json.Unmarshal(thisByte, &outSet.holder)
	return outSet
}
