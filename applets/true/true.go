package true

import (
	"io"
	"os"
)

func Main(stdout io.Writer, args []string) error {
	os.Exit(0)
	return nil
}
