package head

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	opt := &Option{
		bytes: 23,
	}
	f, _ := os.Open("head_test.go")
	w := new(bytes.Buffer)
	head(w, f, opt)
	assert.Equal(t, `package head

import (
`, w.String())
}

func TestLines(t *testing.T) {
	opt := &Option{
		lines: 3,
	}
	f, _ := os.Open("head_test.go")
	w := new(bytes.Buffer)
	head(w, f, opt)
	assert.Equal(t, `package head

import (
`, w.String())
}
