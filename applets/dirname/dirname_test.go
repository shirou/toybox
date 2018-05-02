package dirname

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirname(t *testing.T) {
	var tests = []struct {
		in  []string
		opt *Option
		out string
	}{
		{
			in:  []string{"/tmp"},
			opt: &Option{zeroFlag: true},
			out: "/",
		},
		{
			in:  []string{"/tmp/"},
			opt: &Option{zeroFlag: true},
			out: "/",
		},
		{
			in:  []string{"/tmp/file"},
			opt: &Option{zeroFlag: true},
			out: "/tmp",
		},
	}
	for _, test := range tests {
		w := new(bytes.Buffer)
		assert.Nil(t, dirname(w, test.in, test.opt))
		assert.Equal(t, test.out, w.String())
	}
}

func TestOpt(t *testing.T) {
	{
		w := new(bytes.Buffer)
		opt := &Option{zeroFlag: true}
		assert.Nil(t, dirname(w, []string{"/tmp/file", "/opt/file"}, opt))
		assert.Equal(t, "/tmp/opt", w.String())
	}
	{
		w := new(bytes.Buffer)
		opt := &Option{zeroFlag: false}
		assert.Nil(t, dirname(w, []string{"/tmp/file", "/opt/file"}, opt))
		assert.Equal(t, "/tmp\n/opt\n", w.String())
	}

}
