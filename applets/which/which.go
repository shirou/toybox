package which

import (
	"fmt"
	"io"
	"os/exec"
)

func Main(stdout io.Writer, args []string) error {
	for _, path := range args {
		p, err := exec.LookPath(path)
		if err != nil {
			ne, ok := err.(*exec.Error)
			// which ignore NotFound error
			if ok && ne.Err == exec.ErrNotFound {
				continue
			}
			return err
		}
		fmt.Fprintln(stdout, p)
	}
	return nil
}
