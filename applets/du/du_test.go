package du

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	w := new(bytes.Buffer)
	opt := &Option{depth: 10}

	assert.Nil(t, walk(w, "../", opt))
	fmt.Println(w.String())
}
