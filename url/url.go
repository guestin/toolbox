package url

import (
	"bytes"
	"fmt"
	netUrl "net/url"
)

type (
	_UrlBuildOptions struct {
		Paths   []string
		Queries netUrl.Values
	}
	UrlBuildOption func(op *_UrlBuildOptions)
)

func WithPath(v ...string) UrlBuildOption {
	return func(op *_UrlBuildOptions) {
		op.Paths = append(op.Paths, v...)
	}
}

func WithQuery(k, v string) UrlBuildOption {
	return WithQuery2(k, v, false)
}

func WithQuery2(k, v string, override bool) UrlBuildOption {
	return func(op *_UrlBuildOptions) {
		if override {
			op.Queries.Del(k)
		}
		op.Queries.Add(k, v)
	}
}

func MakeUrlString(base string, opts ...UrlBuildOption) (string, error) {
	url, err := MakeUrl(base, opts...)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func MakeUrl(base string, opts ...UrlBuildOption) (*netUrl.URL, error) {
	url, err := netUrl.Parse(base)
	if err != nil {
		return nil, err
	}
	ubo := _UrlBuildOptions{
		Paths:   []string{},
		Queries: url.Query(),
	}
	for _, op := range opts {
		op(&ubo)
	}
	if len(ubo.Paths) != 0 {
		pathBuilder := bytes.NewBufferString(url.Path)
		for _, pathIt := range ubo.Paths {
			pathBuilder.WriteString(fmt.Sprintf("/%s", pathIt))
		}
		url.Path = pathBuilder.String()
	}
	url.RawQuery = ubo.Queries.Encode()
	return url, nil
}
