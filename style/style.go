package style

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	colorRed    = lipgloss.Color("#FF2D2D")
	colorOrange = lipgloss.Color("#FF8C1A")
	colorGray   = lipgloss.Color("#BDBDBD")
	colorBlue   = lipgloss.Color("#4A90E2")
	colorGreen  = lipgloss.Color("#2ECC71")
	colorGold   = lipgloss.Color("#bda20e")
)

func GetPrimaryColor() color.Color {
	return colorOrange
}

func GetGoldenColor() color.Color {
	return colorGold
}

func GetDangerColor() color.Color {
	return colorRed
}

func GetInfoColor() color.Color {
	return colorBlue
}

func GetSuccessColor() color.Color {
	return colorGreen
}

func GetNeutralColor() color.Color {
	return colorGray
}

func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02X%02X%02X", r>>8, g>>8, b>>8)
}
