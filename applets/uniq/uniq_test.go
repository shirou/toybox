package uniq

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	opt := &Option{
		duplicate: true,
	}

	w := new(bytes.Buffer)
	in := bytes.NewBufferString("AAA\nBBB\n")
	assert.Nil(t, uniq(in, w, opt))

	fmt.Println(w.String())
}
