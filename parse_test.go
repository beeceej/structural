package structural

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func golden(t *testing.T) string {
	f, err := os.Open("testdata/golden")
	assert.NoError(t, err)
	b, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	return string(b)
}

func TestGenerate(t *testing.T) {
	var (
		actual = bytes.NewBuffer(nil)
		expect = golden(t)
	)

	assert.NoError(t, Generate(actual, "testdata/definition.go"))
	assert.Equal(t, expect, actual.String())
}
