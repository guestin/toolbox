package mob

import (
	"testing"
)

func TestSet(t *testing.T) {
	s1 := NewSet()
	SetAdd(s1, 1, 2, 4, 5)
	s2 := NewSet()
	SetAdd(s2, 4, 5, 6, 7, 8)
	for v := range s1.Intersect(s2).Iter() {
		t.Log(v)
	}
}
