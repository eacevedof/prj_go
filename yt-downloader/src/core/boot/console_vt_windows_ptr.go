//go:build windows

package boot

import "unsafe"

func unsafePointer(mode *uint32) unsafe.Pointer {
	return unsafe.Pointer(mode)
}
