//go:build windows

package boot

import (
	"os"
	"syscall"
)

// enableVirtualTerminal el conhost clasico (doble clic en el .exe) no interpreta los codigos ANSI de color hasta que se
// activa ENABLE_VIRTUAL_TERMINAL_PROCESSING; Windows Terminal ya lo trae. Devuelve false si no se pudo (consola
// antigua o salida redirigida): entonces se escribe sin color.
func enableVirtualTerminal() bool {
	const enableVirtualTerminalProcessing = 0x0004
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	enabled := false
	for _, file := range []*os.File{os.Stdout, os.Stderr} {
		handle := file.Fd()
		var mode uint32
		if ok, _, _ := getConsoleMode.Call(handle, uintptr(unsafePointer(&mode))); ok == 0 {
			continue
		}
		if mode&enableVirtualTerminalProcessing != 0 {
			enabled = true
			continue
		}
		if ok, _, _ := setConsoleMode.Call(handle, uintptr(mode|enableVirtualTerminalProcessing)); ok != 0 {
			enabled = true
		}
	}
	return enabled
}
