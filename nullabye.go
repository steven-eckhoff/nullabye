package nullabye

import (
	"encoding/json"
)

var (
	// jsonMarshal is overwritten in tests to force errors
	jsonMarshal = json.Marshal
	// jsonUnmarshal is overwritten in tests to force errors
	jsonUnmarshal = json.Unmarshal
)

type OptionalStruct struct {
	set  bool
	data []byte
}

func NewOptionalStruct(target interface{}) (OptionalStruct, error) {
	data, err := jsonMarshal(target)
	return OptionalStruct{set: err == nil, data: data}, err
}

func NewOptionalStructOrPanic(target interface{}) OptionalStruct {
	os, err := NewOptionalStruct(target)
	if err != nil {
		panic(err)
	}

	return os
}

func (os *OptionalStruct) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		os.set = false
		return nil
	}

	os.set = true
	os.data = data

	return nil
}

func (os OptionalStruct) MarshalJSON() ([]byte, error) {
	if os.set == false {
		return []byte("null"), nil
	}

	return os.data, nil
}

func (os *OptionalStruct) IsSet() bool {
	return os.set
}

func (os *OptionalStruct) Get(target interface{}) (interface{}, error) {
	if os.set == false {
		return target, nil
	}

	err := jsonUnmarshal(os.data, target)
	return target, err
}

func (os *OptionalStruct) GetOrPanic(target interface{}) interface{} {
	target, err := os.Get(target)
	if err != nil {
		panic(err)
	}

	return target
}
