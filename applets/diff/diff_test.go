package diff

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	const root = "testdata"
	var tests = []struct {
		in  string
		opt *Option
	}{
		{in: "a", opt: &Option{unifiedFlag: true}},
	}

	for _, test := range tests {
		files := []string{
			filepath.Join(root, test.in+"-1.txt"),
			filepath.Join(root, test.in+"-2.txt")}

		w := new(bytes.Buffer)
		if err := diff(w, files, test.opt); err != nil {
			assert.Nil(t, err)
			continue
		}
		e := filepath.Join(root, test.in+"-result.txt")
		expected, _ := ioutil.ReadFile(e)

		assert.Equal(t, w.String(), string(expected))
	}
}
