package commands

import (
	"fmt"
	"os"

	"yt-downloader-go/src/modules/shared/infrastructure/traits"
)

// progressPrinter una sola linea por stderr que se sobrescribe: "  45% 12.3 MB / 27.1 MB". Compartido por los
// comandos que descargan; Reset() antes de cada fichero.
type progressPrinter struct {
	traits.CliColorTrait
	lastPercent int
}

func newProgressPrinter() *progressPrinter {
	return &progressPrinter{lastPercent: -1}
}

func (p *progressPrinter) Reset() {
	p.lastPercent = -1
}

func (p *progressPrinter) Print(read int64, total int64) {
	if total <= 0 {
		fmt.Fprintf(os.Stderr, "\r  %s", getHumanBytes(read))
		return
	}
	percent := int(read * 100 / total)
	if percent == p.lastPercent {
		return
	}
	p.lastPercent = percent
	fmt.Fprintf(os.Stderr, "\r  %s", p.GetGreen(fmt.Sprintf("%3d%% %s / %s", percent, getHumanBytes(read), getHumanBytes(total))))
}

func getHumanBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	value := float64(bytes)
	units := []string{"KB", "MB", "GB"}
	i := -1
	for value >= unit && i < len(units)-1 {
		value /= unit
		i++
	}
	return fmt.Sprintf("%.1f %s", value, units[i])
}
