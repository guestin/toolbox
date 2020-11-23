package mob

import mapset "github.com/deckarep/golang-set"

func NewSet() mapset.Set {
	return mapset.NewThreadUnsafeSet()
}

func NewConcurrentSet() mapset.Set {
	return mapset.NewSet()
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
