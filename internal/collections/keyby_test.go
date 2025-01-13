package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyBy(t *testing.T) {
	type aa struct {
		Name string
		Age  int
	}
	type fields struct {
		src []aa
	}
	type args struct {
		key  string
		dest map[int]aa
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test KeyBy",
			fields: fields{
				src: []aa{
					{Name: "John", Age: 20},
					{Name: "JohnA", Age: 22},
				},
			},
			args: args{
				key:  "Age",
				dest: make(map[int]aa),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			KeyBy(tt.args.key, tt.fields.src, tt.args.dest)
			t.Log(tt.args.dest)
			assert.Equal(t, tt.args.dest, map[int]aa{
				20: {Name: "John", Age: 20},
				22: {Name: "JohnA", Age: 22},
			}, "they should be equal")
		})
	}
}

func TestKeyBy2Array(t *testing.T) {
	type aa struct {
		Name string
		Age  int
	}
	type fields struct {
		src []aa
	}
	type args struct {
		key  string
		dest map[int][]aa
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test KeyBy2Array",
			fields: fields{
				src: []aa{
					{Name: "John", Age: 20},
					{Name: "JohnA", Age: 22},
				},
			},
			args: args{
				key:  "Age",
				dest: make(map[int][]aa),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			KeyBy2Array(tt.args.key, tt.fields.src, tt.args.dest)
			t.Log(tt.args.dest)
			assert.Equal(t, tt.args.dest, map[int][]aa{
				20: {{Name: "John", Age: 20}},
				22: {{Name: "JohnA", Age: 22}},
			}, "they should be equal")
		})
	}
}
