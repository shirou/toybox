package false

import (
	"io"
	"os"
)

func Main(stdout io.Writer, args []string) error {
	os.Exit(1)
	return nil
}
