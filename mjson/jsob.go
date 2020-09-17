package mjson

import "encoding/json"

const KJsonIndent = "  "
const KJsonIndentPrefix = ""

type Object map[string]interface{}

type ObjectLessFun func(l Object, r Object) bool

func NewObjectFromBytes(bytes []byte) (Object, error) {
	out := NewObject()
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func NewObject() Object {
	return make(Object)
}

func (this Object) Get(k string) interface{} {
	return this[k]
}

func (this Object) GetOr(k string, def interface{}) interface{} {
	if v, ok := this[k]; ok {
		return v
	}
	return def
}

func (this Object) HasKey(k string) bool {
	_, ok := this[k]
	return ok
}

func (this Object) GetStringOr(k string, def string) string {
	v, ok := this[k]
	if !ok {
		return def
	}
	return any2String(v, def)
}

func (this Object) GetString(k string) string {
	return this.GetStringOr(k, "")
}

func (this Object) GetIntOr(k string, def int) int {
	v, ok := this[k]
	if !ok {
		return def
	}
	return any2Int(v, def)
}

func (this Object) GetBool(k string) bool {
	v, ok := this[k]
	if !ok {
		return false
	}
	if vv, ok := v.(bool); ok {
		return vv
	}
	if vp, ok := v.(*bool); ok {
		return *vp
	}
	return false
}

func (this Object) GetInt(k string) int {
	return this.GetIntOr(k, 0)
}

func (this Object) GetJsonObject(k string) Object {
	v, ok := this[k]
	if !ok {
		return nil
	}
	if rv, ok := v.(map[string]interface{}); ok {
		return rv
	}
	return nil
}

func (this Object) GetJsonArray(k string) Array {
	v, ok := this[k]
	if !ok {
		return nil
	}
	if rv, ok := v.([]interface{}); ok {
		return rv
	}
	return nil
}

func (this Object) ToJsonString() (string, error) {
	return MarshalObject(this)
}

// 注意,在函数之间传递时,必须使用指针进行传递,
// 否则后续任何的修改都是行为未定义
type Array []interface{}

const DefaultJsonArrayCap = 4

func NewArray() Array {
	return NewArrayWithCap(DefaultJsonArrayCap)
}

func NewArrayWithCap(cap int) Array {
	ret := make(Array, 0, cap)
	return ret
}

func (this *Array) Add(o interface{}) *Array {
	*this = append(*this, o)
	return this
}

func (this Array) ToJsonString() (string, error) {
	return MarshalObject(this)
}

//将数据转换为json
func MarshalObject(e interface{}) (string, error) {
	if bs, err := json.Marshal(e); err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}
