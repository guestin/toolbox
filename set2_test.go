package mob

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestSet(t *testing.T) {
	x := []interface{}{1, 2, 4, 5}
	s1 := NewSetFromSlice(x)
	SetAdd(s1, 1, 2, 4, 5)
	s2 := NewSet()
	SetAdd(s2, 4, 5, 6, 7, 8)
	t.Log("s1 diff s2", s1.Difference(s2))
	t.Log("s2 diff s1", s2.Difference(s1))
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
