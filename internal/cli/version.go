package cli

import "fmt"

const appVersion = "0.1.0"

func (a *App) runVersion() int {
	fmt.Printf("compver version %s\n", appVersion)
	return 0
}
