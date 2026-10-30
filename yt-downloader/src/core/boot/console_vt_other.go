//go:build !windows

package boot

// enableVirtualTerminal fuera de Windows las terminales ya entienden ANSI.
func enableVirtualTerminal() bool {
	return true
}
