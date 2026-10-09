package cli

import "fmt"

func (a *App) runHelp() int {
	fmt.Println("Uso:")
	fmt.Println("  compver [ruta]")
	fmt.Println("  compver scan [ruta]")
	fmt.Println("  compver version")
	fmt.Println("  compver help")
	fmt.Println()
	fmt.Println("Ejemplos:")
	fmt.Println("  compver scan .")
	fmt.Println("  compver scan /home/NAME/proyectos/mi-app")
	fmt.Println("  compver scan --tool go .")
	fmt.Println("  compver scan --tool terraform .")
	fmt.Println("  compver version")

	return 0
}
