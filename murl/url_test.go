package murl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeUrl(t *testing.T) {
	res, err := MakeUrlString("https://api.guesin.cn/a/b?query=1",
		WithQuery("name", "中文"),
		WithQuery("what?", "123"),
		WithPath("c", "d", "e"),
	)
	assert.NoError(t, err)
	t.Log(res)
}
