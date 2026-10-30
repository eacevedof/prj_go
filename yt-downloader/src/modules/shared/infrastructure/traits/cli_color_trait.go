package traits

import (
	"fmt"

	"yt-downloader-go/src/core/boot"
)

// CliColorTrait colores ANSI para la consola: GetXxx() devuelve el texto coloreado, EchoXxx() lo imprime por stdout.
// Si la terminal no los soporta (boot.AreColorsEnabled), devuelve el texto tal cual.
type CliColorTrait struct{}

func (t CliColorTrait) paint(code string, text string) string {
	if !boot.AreColorsEnabled() {
		return text
	}
	return code + text + "\033[0m"
}

func (t CliColorTrait) EchoRed(text string) {
	fmt.Println(t.GetRed(text))
}

func (t CliColorTrait) GetRed(text string) string {
	return t.paint("\033[91m", text)
}

func (t CliColorTrait) EchoGreen(text string) {
	fmt.Println(t.GetGreen(text))
}

func (t CliColorTrait) GetGreen(text string) string {
	return t.paint("\033[92m", text)
}

func (t CliColorTrait) EchoYellow(text string) {
	fmt.Println(t.GetYellow(text))
}

func (t CliColorTrait) GetYellow(text string) string {
	return t.paint("\033[93m", text)
}

func (t CliColorTrait) EchoBlue(text string) {
	fmt.Println(t.GetBlue(text))
}

func (t CliColorTrait) GetBlue(text string) string {
	return t.paint("\033[94m", text)
}

func (t CliColorTrait) EchoWhite(text string) {
	fmt.Println(t.GetWhite(text))
}

func (t CliColorTrait) GetWhite(text string) string {
	return t.paint("\033[97m", text)
}
