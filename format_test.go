package spec

import (
	"gotest.tools/assert"
	"testing"
)

func Test_Format_Integer_Pass(t *testing.T) {
	err := Integer.Check("123")
	assert.Equal(t, err == nil, true)

	err = Integer.Check("+123")
	assert.Equal(t, err == nil, true)

	err = Integer.Check("-123")
	assert.Equal(t, err == nil, true)

	err = Integer.Check("0")
	assert.Equal(t, err == nil, true)
}

func Test_Format_Integer_Fail(t *testing.T) {
	err := Integer.Check("123.4")
	assert.Equal(t, err != nil, true)

	err = Integer.Check("abcd")
	assert.Equal(t, err != nil, true)

	err = Integer.Check("-")
	assert.Equal(t, err != nil, true)

	err = Integer.Check("+")
	assert.Equal(t, err != nil, true)
}
