package mob

import mapset "github.com/deckarep/golang-set"

func NewSet() mapset.Set {
	return mapset.NewThreadUnsafeSet()
}

func NewSetFrom(vs ...interface{}) mapset.Set {
	return NewSetFromSlice(vs)
}

func NewSetFromSlice(vs []interface{}) mapset.Set {
	return mapset.NewThreadUnsafeSetFromSlice(vs)
}

func NewConcurrentSet() mapset.Set {
	return mapset.NewSet()
}

func NewConcurrentSetFrom(vs ...interface{}) mapset.Set {
	return NewConcurrentSetFromSlice(vs)
}

func NewConcurrentSetFromSlice(vs ...interface{}) mapset.Set {
	return mapset.NewSetFromSlice(vs)
}

func SetAdd(dst mapset.Set, vs ...interface{}) mapset.Set {
	return SetFromSlice(dst, vs)
}

func SetFromSlice(dst mapset.Set, vs []interface{}) mapset.Set {
	for _, it := range vs {
		dst.Add(it)
	}
	return dst
}
