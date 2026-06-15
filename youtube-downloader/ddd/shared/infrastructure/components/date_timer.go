package components

import "time"

// DateTimer provides date/time formatting helpers.
type DateTimer struct{}

// NewDateTimer returns a ready-to-use DateTimer.
func NewDateTimer() *DateTimer { return &DateTimer{} }

func (d *DateTimer) GetCurrentDateYmd() string {
	return time.Now().Format("2006-01-02")
}

func (d *DateTimer) GetCurrentDatetimeYmdHms() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (d *DateTimer) GetTimestamp() int64 {
	return time.Now().Unix()
}

func (d *DateTimer) GetDatetimeAsYmd(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func (d *DateTimer) GetDatetimeAsYmdHms(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
