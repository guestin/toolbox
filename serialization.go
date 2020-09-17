package mod

import (
	"encoding/json"
)

type Serialization interface {
	Marshal(v interface{}) ([]byte, error)
	UnMarshal(data []byte, v interface{}) error
}

type binderProxy struct {
	marshalFn   func(v interface{}) ([]byte, error)
	unmarshalFn func(data []byte, v interface{}) error
}

func (this *binderProxy) Marshal(v interface{}) ([]byte, error) {
	return this.marshalFn(v)
}

func (this *binderProxy) UnMarshal(data []byte, v interface{}) error {
	return this.unmarshalFn(data, v)
}

var jsonSerialization Serialization = &binderProxy{
	marshalFn:   json.Marshal,
	unmarshalFn: json.Unmarshal,
}

func JsonSerialization() Serialization {
	return jsonSerialization
}
