package head

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	out := new(bytes.Buffer)
	out.Grow(25)
	assert.Nil(t, head(out, OpenRequiredFile("testdata/loremipsum.txt", t), &Option{bytes: 23}))
	assert.Equal(t, testBytesExpectedResult, out.String())
}

func TestLines(t *testing.T) {
	out := new(bytes.Buffer)
	out.Grow(250)
	assert.Nil(t, head(out, OpenRequiredFile("testdata/loremipsum.txt", t), &Option{lines: 3}))
	assert.Equal(t, testLinesExpectedResult, out.String())
}

func TestHeadNegativeLines(t *testing.T) {
	out := new(bytes.Buffer)
	out.Grow(550)
	assert.Nil(t, head(out, OpenRequiredFile("testdata/loremipsum.txt", t), &Option{lines: -3}))
	assert.Equal(t, testHeadNegativeLinesExpectedResult, out.String())
}

func TestHeadNegativeBytes(t *testing.T) {
	out := new(bytes.Buffer)
	out.Grow(600)
	assert.Nil(t, head(out, OpenRequiredFile("testdata/loremipsum.txt", t), &Option{bytes: -175}))
	assert.Equal(t, testHeadNegativeBytesExpectedResult, out.String())
}

func OpenRequiredFile(path string, t *testing.T) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("cannot open required file '%s': %s", path, err.Error())
	}
	return f
}

const testBytesExpectedResult = "Lorem ipsum dolor sit a"
const testLinesExpectedResult = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis eleifend, velit
sed pulvinar tincidunt, nulla ante gravida massa, non vulputate mauris ex ac
velit. Mauris vehicula ipsum ut lobortis tristique. Curabitur vel neque et eros
`
const testHeadNegativeLinesExpectedResult = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis eleifend, velit
sed pulvinar tincidunt, nulla ante gravida massa, non vulputate mauris ex ac
velit. Mauris vehicula ipsum ut lobortis tristique. Curabitur vel neque et eros
dapibus vehicula. Nulla lectus tellus, molestie ut cursus at, feugiat gravida
risus. Nulla tincidunt est sed arcu tristique, eu convallis massa sodales.
Pellentesque ornare eros nec ex elementum feugiat. Nulla condimentum lectus
ipsum, id lacinia ex feugiat eu. Vivamus arcu arcu, sollicitudin id aliquam ut,
`
const testHeadNegativeBytesExpectedResult = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis eleifend, velit
sed pulvinar tincidunt, nulla ante gravida massa, non vulputate mauris ex ac
velit. Mauris vehicula ipsum ut lobortis tristique. Curabitur vel neque et eros
dapibus vehicula. Nulla lectus tellus, molestie ut cursus at, feugiat gravida
risus. Nulla tincidunt est sed arcu tristique, eu convallis massa sodales.
Pellentesque ornare eros nec ex elementum feugiat. Nulla condimentum lectus
ipsum, id lacinia ex feugiat eu. Vivamus arcu arcu, sollicitudin id aliquam ut,
luctus id erat. Nunc diam`
