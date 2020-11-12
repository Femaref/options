package options

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInline(t *testing.T) {
	expt := New().Expect("foo", reflect.TypeOf(""), "", false)

	in := Options{
		"foo": "bar",
	}

	var out Options

	err := expt.Parse(&out, in)

	assert.NoError(t, err)
	assert.Equal(t, in["foo"], out["foo"])
}

func TestStruct(t *testing.T) {
	expt := New().Expect("foo", reflect.TypeOf(""), "", false)

	in := Options{
		"foo": "bar",
	}

	var out struct {
		Foo string
	}

	err := expt.Parse(&out, in)

	assert.NoError(t, err)
	assert.Equal(t, in["foo"], out.Foo)
}
