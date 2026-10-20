package cli

import (
	"fmt"
	"os"

	"github.com/S4M73l09/compver/internal/model"
)

type terminalStyle struct {
	enabled bool
}

func newTerminalStyle(noColor bool) terminalStyle {
	return terminalStyle{
		enabled: !noColor && isTerminal(os.Stdout),
	}
}

func (s terminalStyle) heading(value string) string {
	return s.paint("\033[1;36m", value)
}

func (s terminalStyle) status(value string) string {
	icon := "?"
	color := "\033[90m"

	switch value {
	case string(model.StatusCurrent):
		icon = "✓"
		color = "\033[32m"
	case string(model.StatusUpdateAvailable), string(model.StatusCompatibleUpdate):
		icon = "↑"
		color = "\033[33m"
	case string(model.StatusVulnerable), string(model.StatusVulnerableUpdate):
		icon = "⚠"
		color = "\033[31m"
	case string(model.StatusEndOfLife):
		icon = "✖"
		color = "\033[1;31m"
	case string(model.StatusDeprecated):
		icon = "!"
		color = "\033[33m"
	}

	return s.paint(color, icon+" "+value)
}

func (s terminalStyle) count(value int) string {
	color := "\033[32m"
	if value > 0 {
		color = "\033[33m"
	}
	return s.paint(color, fmt.Sprintf("%d", value))
}

func (s terminalStyle) paint(color string, value string) string {
	if !s.enabled {
		return value
	}
	return color + value + "\033[0m"
}
