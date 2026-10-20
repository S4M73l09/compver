package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/S4M73l09/compver/internal/providers"
)

type scanOptions struct {
	Tool               string
	IncludePreReleases bool
	AllVersions        bool
	Limit              int
	NetworkMode        providers.NetworkMode
	Show               map[string]bool
	OutputFormat       string
	NoColor            bool
	Verbose            bool
	Summary            bool
	Quiet              bool
	FailOn             map[string]bool
}

const (
	showLifecycle       = "lifecycle"
	showVulnerabilities = "vulnerabilities"
	showSources         = "sources"
	showReleases        = "releases"
)

func parseScanOptions(args []string) (
	scanOptions,
	string,
	error,
) {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	normalizedArgs, err := normalizeScanArgs(args)
	if err != nil {
		return scanOptions{}, "", err
	}

	tool := flags.String(
		"tool",
		"",
		"herramienta que se desea analizar",
	)

	includePreReleases := flags.Bool(
		"include-prereleases",
		false,
		"incluir alpha, beta y rc",
	)

	allVersions := flags.Bool(
		"all-versions",
		false,
		"mostrar todas las versiones",
	)

	limit := flags.Int(
		"limit",
		3,
		"número máximo de versiones",
	)

	offline := flags.Bool(
		"offline",
		false,
		"no realizar consultas de red",
	)

	refresh := flags.Bool(
		"refresh",
		false,
		"ignorar la caché y consultar de nuevo",
	)

	show := flags.String(
		"show",
		"",
		"secciones adicionales: lifecycle, vulnerabilities, sources, releases",
	)

	format := flags.String(
		"format",
		"table",
		"formato de salida: table o json",
	)

	noColor := flags.Bool(
		"no-color",
		false,
		"desactivar colores de la terminal",
	)

	verbose := flags.Bool(
		"verbose",
		false,
		"mostrar todos los detalles del análisis",
	)

	summary := flags.Bool(
		"summary",
		false,
		"mostrar únicamente los resúmenes",
	)

	quiet := flags.Bool(
		"quiet",
		false,
		"silenciar la salida normal",
	)

	failOn := flags.String(
		"fail-on",
		"",
		"fallar si se detectan: vulnerable, end-of-life",
	)

	if err := flags.Parse(normalizedArgs); err != nil {
		return scanOptions{}, "", err
	}

	remainingArgs := flags.Args()

	if len(remainingArgs) > 1 {
		return scanOptions{}, "", fmt.Errorf(
			"scan acepta como máximo una ruta",
		)
	}

	if *limit < 1 && !*allVersions {
		return scanOptions{}, "", fmt.Errorf(
			"limit debe ser mayor que cero",
		)
	}

	if *offline && *refresh {
		return scanOptions{}, "", fmt.Errorf(
			"offline y refresh no pueden utilizarse juntos",
		)
	}

	presentationModes := 0
	if *verbose {
		presentationModes++
	}
	if *summary {
		presentationModes++
	}
	if *quiet {
		presentationModes++
	}
	if presentationModes > 1 {
		return scanOptions{}, "", fmt.Errorf(
			"verbose, summary y quiet no pueden utilizarse juntos",
		)
	}

	if *format != "table" && *format != "json" {
		return scanOptions{}, "", fmt.Errorf(
			"formato no soportado: %s (válidos: table, json)",
			*format,
		)
	}

	networkMode := providers.NetworkAuto
	if *offline {
		networkMode = providers.NetworkOffline
	}
	if *refresh {
		networkMode = providers.NetworkRefresh
	}

	if *tool != "" && *tool != "go" && *tool != "terraform" {
		return scanOptions{}, "", fmt.Errorf(
			"herramienta no soportada: %s",
			*tool,
		)
	}

	showSections, err := parseShowSections(*show)
	if err != nil {
		return scanOptions{}, "", err
	}
	if *verbose {
		showSections[showLifecycle] = true
		showSections[showVulnerabilities] = true
		showSections[showSources] = true
		showSections[showReleases] = true
	}

	failConditions, err := parseFailOn(*failOn)
	if err != nil {
		return scanOptions{}, "", err
	}

	root := "."

	if len(remainingArgs) == 1 {
		root = remainingArgs[0]
	}

	return scanOptions{
		Tool:               *tool,
		IncludePreReleases: *includePreReleases,
		AllVersions:        *allVersions,
		Limit:              *limit,
		NetworkMode:        networkMode,
		Show:               showSections,
		OutputFormat:       *format,
		NoColor:            *noColor,
		Verbose:            *verbose,
		Summary:            *summary,
		Quiet:              *quiet,
		FailOn:             failConditions,
	}, root, nil
}

func normalizeScanArgs(args []string) ([]string, error) {
	valueFlags := map[string]bool{
		"--tool":    true,
		"--limit":   true,
		"--show":    true,
		"--format":  true,
		"--fail-on": true,
	}

	flags := make([]string, 0, len(args))
	positional := make([]string, 0, 1)

	for index := 0; index < len(args); index++ {
		argument := args[index]

		if argument == "--" {
			positional = append(positional, args[index+1:]...)
			break
		}

		if !strings.HasPrefix(argument, "-") {
			positional = append(positional, argument)
			continue
		}

		flags = append(flags, argument)

		flagName := argument
		if equalsIndex := strings.IndexByte(argument, '='); equalsIndex >= 0 {
			flagName = argument[:equalsIndex]
		}

		if !valueFlags[flagName] || strings.Contains(argument, "=") {
			continue
		}

		if index+1 >= len(args) {
			return nil, fmt.Errorf("la opción %s requiere un valor", argument)
		}

		index++
		flags = append(flags, args[index])
	}

	return append(flags, positional...), nil
}

func parseShowSections(value string) (map[string]bool, error) {
	sections := make(map[string]bool)
	if strings.TrimSpace(value) == "" {
		return sections, nil
	}

	validSections := map[string]bool{
		showLifecycle:       true,
		showVulnerabilities: true,
		showSources:         true,
		showReleases:        true,
	}

	for _, rawSection := range strings.Split(value, ",") {
		section := strings.TrimSpace(strings.ToLower(rawSection))
		if section == "" {
			continue
		}
		if !validSections[section] {
			return nil, fmt.Errorf(
				"sección no soportada en --show: %s (válidas: lifecycle, vulnerabilities, sources, releases)",
				section,
			)
		}
		sections[section] = true
	}

	return sections, nil
}
