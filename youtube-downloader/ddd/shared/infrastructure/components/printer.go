package components

import (
	"fmt"
	"os"
)

const (
	printerColorRed    = "\033[31m"
	printerColorGreen  = "\033[32m"
	printerColorYellow = "\033[33m"
	printerColorBlue   = "\033[34m"
	printerColorOrange = "\033[38;5;208m"
	printerColorWhite  = "\033[37m"
	printerColorLemon  = "\033[93m"
	printerColorReset  = "\033[0m"
)

// Printer writes ANSI-colored lines to stdout.
type Printer struct{}

// NewPrinter returns a ready-to-use Printer.
func NewPrinter() *Printer { return &Printer{} }

func (p *Printer) PrintRed(text string) string    { return p.colorize(printerColorRed, text) }
func (p *Printer) PrintGreen(text string) string  { return p.colorize(printerColorGreen, text) }
func (p *Printer) PrintYellow(text string) string { return p.colorize(printerColorYellow, text) }
func (p *Printer) PrintBlue(text string) string   { return p.colorize(printerColorBlue, text) }
func (p *Printer) PrintOrange(text string) string { return p.colorize(printerColorOrange, text) }
func (p *Printer) PrintWhite(text string) string  { return p.colorize(printerColorWhite, text) }
func (p *Printer) PrintLemon(text string) string  { return p.colorize(printerColorLemon, text) }

// Die prints text in red then exits with code 1.
func (p *Printer) Die(text string) {
	p.PrintRed(text)
	os.Exit(1)
}

func (p *Printer) colorize(color, text string) string {
	colored := color + text + printerColorReset
	fmt.Println(colored)
	return colored
}
