package tr

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTr(t *testing.T) {
	w := new(bytes.Buffer)
	r := new(bytes.Buffer)
	fmt.Fprint(r, "111")

	tr(w, r, "1", "22")
	assert.Equal(t, "222222\n", w.String())
}
