package uuidgen

import (
	"fmt"
	"io"
	"syscall"
	"unsafe"
)

var (
	rpcrt4Dll      = syscall.NewLazyDLL("rpcrt4.dll")
	procUuidCreate = rpcrt4Dll.NewProc("UuidCreate")
)

const RPC_S_OK = 0

func uuidgen(w io.Writer) error {
	var uuid [16]byte
	rc, _, e := syscall.Syscall(procUuidCreate.Addr(), 1,
		uintptr(unsafe.Pointer(&uuid[0])), 0, 0)
	if int(rc) != RPC_S_OK {
		if e != 0 {
			return error(e)
		} else {
			return nil
		}
	}
	fmt.Fprintln(w, string(uuid[:]))
	return nil
}
