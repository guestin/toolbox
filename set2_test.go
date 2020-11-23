package mob

import (
	"github.com/stretchr/testify/assert"
	"net/url"
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

func TestURL(t *testing.T) {
	u, err := url.Parse("https://foo.com/a/b?q=100")
	assert.NoError(t, err)
	t.Log(u.Path)
	queries := u.Query()
	queries.Add("name", "中文")
	u.RawQuery = queries.Encode()
	t.Log(u.String())
}
