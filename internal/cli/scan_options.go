package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/S4M73l09/compver/internal/providers"
)

type scanOptions struct {
	Tool               string
	IncludePreReleases bool
	AllVersions        bool
	Limit              int
	NetworkMode        providers.NetworkMode
}

func parseScanOptions(args []string) (
	scanOptions,
	string,
	error,
) {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

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

	if err := flags.Parse(args); err != nil {
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
	}, root, nil
}
