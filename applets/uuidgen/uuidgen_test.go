package uuidgen

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	w := new(bytes.Buffer)
	uuidgen(w)
	assert.Equal(t, 5, len(strings.Split(w.String(), "-")))
}
