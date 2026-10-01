package cli

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		return a.runScan([]string{"."})
	}

	switch args[0] {
	case "scan":
		return a.runScan(args[1:])

	case "version":
		return a.runVersion()

	case "help", "-h", "--help":
		return a.runHelp()

	default:
		return a.runScan(args)
	}
}
