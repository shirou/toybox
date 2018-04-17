package cut

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	var tests = []struct {
		in  string
		opt Option
		out string
	}{
		{"abcd", Option{charPos: "3"}, "c\n"},                          // a-character
		{"abcd", Option{charPos: "1-2"}, "ab\n"},                       // a-closed-range
		{"f1\tf2\tf3", Option{fieldPos: "2", delimiter: "\t"}, "f2\n"}, // a-field
		{"abcd", Option{charPos: "-3"}, "abc\n"},                       // an-open-range
		{"abcd", Option{charPos: "3-"}, "cd\n"},                        // an-unclosed-range
		{`over
quick`,
			Option{charPos: "3-"}, `er
ick
`},
	}

	for _, test := range tests {
		w := new(bytes.Buffer)
		r := bytes.NewBufferString(test.in)

		assert.Nil(t, cut(w, r, test.opt))
		assert.Equal(t, test.out, w.String())
	}

}
