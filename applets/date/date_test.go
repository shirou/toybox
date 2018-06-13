package date

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReference(t *testing.T) {
	w := new(bytes.Buffer)

	opt := &Option{
		referenceFlag: true,
	}
	assert.NotNil(t, reference(w, []string{}, opt))
	assert.NotNil(t, reference(w, []string{"date.go", "date.go"}, opt))
	assert.Nil(t, reference(w, []string{"date.go"}, opt))
}

func TestRFC3339(t *testing.T) {
	var tests = []struct {
		opt *Option
		out string
	}{
		{&Option{rfc3339Flag: "date"}, "2018-01-01\n"},
		{&Option{rfc3339Flag: "seconds"}, "2018-01-01 12:34:56Z\n"},
		{&Option{rfc3339Flag: "ns"}, "2018-01-01 12:34:56Z\n"},
		{&Option{iso8601Flag: "date"}, "2018-01-01\n"},
		{&Option{iso8601Flag: "hours"}, "2018-01-01T12Z\n"},
		{&Option{iso8601Flag: "minutes"}, "2018-01-01T12:34Z\n"},
		{&Option{iso8601Flag: "seconds"}, "2018-01-01T12:34:56Z\n"},
		{&Option{iso8601Flag: "ns"}, "2018-01-01T12:34:56Z\n"},
	}

	ti, _ := time.Parse(time.RFC3339, "2018-01-01T12:34:56Z")
	for _, test := range tests {
		w := new(bytes.Buffer)
		assert.Nil(t, output(w, ti, test.opt))
		assert.Equal(t, test.out, w.String())
	}
}
