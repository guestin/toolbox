package merrors

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Error interface {
	SetCode(code int) Error
	SetMsg(msg string) Error
	SetData(data interface{}) Error
	GetData() interface{}
	GetMsg() string
	GetCode() int
	Is(err error) bool
	Error() string
}

type _stdError struct {
	Code    int         `json:"code"`           // 错误码
	Msg     string      `json:"message"`        // 错误信息
	Data    interface{} `json:"data,omitempty"` // data
	inherit error       // 包装的外部error
}

func (this *_stdError) GetCode() int {
	return this.Code
}

func (this *_stdError) GetMsg() string {
	return this.Msg
}

func (this *_stdError) GetData() interface{} {
	return this.Data
}

func (this *_stdError) Error() string {
	return this.Msg
}

func (this *_stdError) Is(err error) bool {
	if err == nil {
		return false
	}
	if err == error(this) || this.inherit == err {
		return true
	}
	if this.inherit == error(this) && this.inherit != err {
		return false
	}
	xe, ok := this.inherit.(*_stdError)
	if ok {
		return xe.Is(err)
	}
	return false
}

func NewError() Error {
	return &_stdError{
		Code:    0,
		Msg:     "success",
		Data:    nil,
		inherit: nil,
	}
}

func newErrorWrapped(err error) *_stdError {
	return &_stdError{
		inherit: err,
	}
}

func Errorf(format string, args ...interface{}) Error {
	return NewError().SetMsg(fmt.Sprintf(format, args...))
}

func ErrorWrap(err error, msg string) Error {
	return newErrorWrapped(err).SetMsg(fmt.Sprintf("%s: %s", msg, err.Error()))
}

func ErrorWrapf(err error, format string, args ...interface{}) Error {
	return ErrorWrap(err, fmt.Sprintf(format, args...))
}

func (this *_stdError) SetMsg(msg string) Error {
	this.Msg = msg
	return this
}

func (this *_stdError) SetCode(code int) Error {
	this.Code = code
	return this
}

func (this *_stdError) SetData(data interface{}) Error {
	this.Data = data
	return this
}

func IsStdError(e error) (err Error, ok bool) {
	if e == nil {
		return nil, false
	}
	err, ok = e.(*_stdError)
	return
}

func AssertError(err error, msg string) {
	if err == nil {
		return
	}
	//defer func() {
	//	Assert(false, output)
	//}()
	panic(errors.Wrapf(err, "%s :causedBy", msg))
}

func Assert(cond bool, msg string) {
	if cond {
		return
	}
	//defer func() {
	//	const abortCode = 6
	//	os.Exit(abortCode)
	//}()
	e := errors.Errorf("assertFailed!,%s", msg)
	panic(e)

}

type CombinedErrors []error

func (this CombinedErrors) Error() string {
	builder := strings.Builder{}
	for idx, it := range this {
		builder.WriteString(fmt.Sprintf("err %d:%s", idx+1, it.Error()))
	}
	return builder.String()
}
