package cli

import (
	"fmt"
	"os"

	"github.com/S4M73l09/compver/internal/app"
	"github.com/S4M73l09/compver/internal/providers"
)

func (a *App) runScan(args []string) int {
	options, root, err := parseScanOptions(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	analyzer, err := app.NewAnalyzer(options.Tool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	spinner := NewSpinner(
		os.Stderr,
		"Analizando dependencias y consultando versiones...",
		isTerminal(os.Stderr),
	)

	spinner.Start()

	result, err := analyzer.AnalyzeWithOptions(
		root,
		providers.QueryOptions{
			Mode: options.NetworkMode,
		},
	)

	spinner.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	policyViolated := scanPolicyViolated(result, options.FailOn)

	if options.OutputFormat == "json" {
		if err := printScanJSON(os.Stdout, result, options); err != nil {
			fmt.Fprintf(os.Stderr, "error generando JSON: %v\n", err)
			return 1
		}

		if policyViolated {
			return 2
		}

		return 0
	}

	if options.Quiet {
		if policyViolated {
			return 2
		}

		return 0
	}

	style := newTerminalStyle(options.NoColor)
	printScanSummary(result, style)

	if len(result.Dependencies) == 0 {
		if policyViolated {
			return 2
		}

		return 0
	}
	if options.Summary {
		printToolSummaries(result, style)
		return 0
	}

	printScanResult(result, options, style)
	if policyViolated {
		return 2
	}

	return 0
}
