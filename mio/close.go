package mio

import "io"

func CloseIgnoreErr(clo io.Closer) {
	if clo == nil {
		return
	}
	_ = clo.Close()
}
