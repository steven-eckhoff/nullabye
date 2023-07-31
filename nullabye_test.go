package nullabye

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOptionalStruct(t *testing.T) {
	type option struct {
		Id     int
		Name   string
		Amount float64
	}
	type inputs struct {
		target option
		err    error
	}
	type expectations struct {
		os  OptionalStruct
		err error
	}
	type testCase struct {
		name  string
		build func(t *testing.T) (inputs, expectations)
	}

	testCases := []testCase{
		{
			name: "error",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						target: option{Id: 1, Name: "test", Amount: 9.99},
						err:    assert.AnError,
					}, expectations{
						os:  OptionalStruct{},
						err: assert.AnError,
					}
			},
		},
		{
			name: "happy",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						target: option{Id: 1, Name: "test", Amount: 9.99},
					}, expectations{
						os: OptionalStruct{set: true, data: []byte(`{"Id":1,"Name":"test","Amount":9.99}`)},
					}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input, expected := tc.build(t)
			jsonMarshal = func(v interface{}) ([]byte, error) {
				if input.err != nil {
					return nil, input.err
				}
				return json.Marshal(v)
			}

			got, err := NewOptionalStruct(input.target)
			assert.Equal(t, expected.err, err)
			assert.Equal(t, expected.os, got)
		})
	}
}

func TestOptionalStruct_UnmarshalJSON(t *testing.T) {
	type option struct {
		Id     int
		Name   string
		Amount float64
	}
	type jsonModel struct {
		Id     int
		Name   string
		Option OptionalStruct
	}

	type inputs struct {
		data []byte
	}
	type expectations struct {
		model     jsonModel
		optionSet bool
		option    *option
	}
	type testCase struct {
		name  string
		build func(t *testing.T) (inputs, expectations)
	}

	testCases := []testCase{
		{
			name: "option_present",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						data: []byte(`{"id": 1, "name": "test", "option": {"id": 2, "name": "optional", "amount": 9.99}}`),
					}, expectations{
						model: jsonModel{
							Id:     1,
							Name:   "test",
							Option: OptionalStruct{set: true, data: []byte(`{"id": 2, "name": "optional", "amount": 9.99}`)},
						},
					}
			},
		},
		{
			name: "option_null",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						data: []byte(`{"id": 1, "name": "test", "option": null}`),
					}, expectations{
						model: jsonModel{Id: 1, Name: "test", Option: OptionalStruct{set: false}},
					}
			},
		},
		{
			name: "option_not_present",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						data: []byte(`{"id": 1, "name": "test"}`),
					}, expectations{
						model: jsonModel{Id: 1, Name: "test", Option: OptionalStruct{set: false}},
					}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input, expected := tc.build(t)
			var got jsonModel
			err := json.Unmarshal(input.data, &got)
			require.NoError(t, err)
			assert.Equal(t, expected.model, got)
		})
	}
}

func TestOptionalStruct_MarshalJSON(t *testing.T) {
	type input struct {
		os OptionalStruct
	}
	type expectations struct {
		data []byte
		err  error
	}
	type testCase struct {
		name  string
		build func(t *testing.T) (input, expectations)
	}

	testCases := []testCase{
		{
			name: "not_set",
			build: func(t *testing.T) (input, expectations) {
				return input{
						os: OptionalStruct{set: false},
					}, expectations{
						data: []byte("null"),
						err:  nil,
					}
			},
		},
		{
			name: "set",
			build: func(t *testing.T) (input, expectations) {
				return input{
						os: OptionalStruct{set: true, data: []byte(`{"id": 2, "name": "optional", "amount": 9.99}`)},
					}, expectations{
						data: []byte(`{"id":2,"name":"optional","amount":9.99}`),
						err:  nil,
					}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input, expected := tc.build(t)
			got, err := json.Marshal(input.os)
			assert.Equal(t, expected.err, err)
			assert.Equal(t, expected.data, got)
		})
	}
}

func TestOptionalStruct_Get(t *testing.T) {
	type option struct {
		Id     int
		Name   string
		Amount float64
	}

	type inputs struct {
		os  OptionalStruct
		err error
	}
	type expectations struct {
		option option
		err    error
	}
	type testCase struct {
		name  string
		build func(t *testing.T) (inputs, expectations)
	}

	testCases := []testCase{
		{
			name: "error",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						os:  OptionalStruct{set: true, data: []byte(`0xdeadbeef`)},
						err: assert.AnError,
					}, expectations{
						option: option{},
						err:    assert.AnError,
					}
			},
		},
		{
			name: "not_set",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						os:  OptionalStruct{set: false},
						err: nil,
					}, expectations{
						option: option{},
					}
			},
		},
		{
			name: "set",
			build: func(t *testing.T) (inputs, expectations) {
				return inputs{
						os:  OptionalStruct{set: true, data: []byte(`{"id": 2, "name": "optional", "amount": 9.99}`)},
						err: nil,
					}, expectations{
						option: option{Id: 2, Name: "optional", Amount: 9.99},
					}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			input, expected := tc.build(t)
			jsonUnmarshal = func(data []byte, v interface{}) error {
				if input.err != nil {
					return input.err
				}
				return json.Unmarshal(data, v)
			}
			var got option
			var gotInterface interface{}
			gotInterface, err := input.os.Get(&got)
			assert.ErrorIs(t, err, expected.err)
			assert.Equal(t, expected.option, got)
			assert.Equal(t, &expected.option, gotInterface)
		})
	}
}
