package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "error running command", "error", err)
	}
}

func run(ctx context.Context) error {
	cmd := &cli.Command{
		Name:      "helm-render-ui",
		ArgsUsage: "[helm chart folder]",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:      "helm-chart-folder",
				UsageText: "A folder containing a Chart.yaml file",
				Min:       1,
				Max:       1,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "namespace",
				Aliases: []string{"n"},
				Usage:   "namespace",
				Value:   "default",
			},
			&cli.StringFlag{
				Name:    "release",
				Aliases: []string{"r"},
				Usage:   "release name",
			},
			&cli.StringSliceFlag{
				Name:    "value-file",
				Aliases: []string{"f"},
				Usage:   "extra configuration values file name",
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			chartFolder := command.StringArgs("helm-chart-folder")[0]

			if strings.TrimSpace(chartFolder) == "" {
				return fmt.Errorf("helm chart folder is required")
			}

			chart, err := loader.LoadDir(chartFolder)
			if err != nil {
				return fmt.Errorf("error loading chart from folder: %w", err)
			}

			var displayValueFiles []string

			values := map[string]any{}
			for _, valueFile := range command.StringSlice("value-file") {
				currentMap := map[string]interface{}{}

				displayValueFiles = append(displayValueFiles, ensureRelativePath(strings.TrimPrefix(valueFile, chartFolder)))

				bytes, err := os.ReadFile(valueFile)
				if err != nil {
					return err
				}

				if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
					return fmt.Errorf("failed to parse %s: %w", valueFile, err)
				}
				// Merge with the previous map
				values = mergeMaps(values, currentMap)
			}

			if err := chartutil.ProcessDependencies(chart, values); err != nil {
				return err
			}

			options := chartutil.ReleaseOptions{
				Name:      command.String("release"),
				Namespace: command.String("namespace"),
				Revision:  1,
				IsInstall: true,
				IsUpgrade: false,
			}
			if options.Name == "" {
				options.Name = chart.Metadata.Name
			}

			browserURL := "http://127.0.0.1:" + httpPort
			slog.InfoContext(ctx, "opening browser URL", "url", browserURL)
			// _ = openURL(browserURL)

			return runHTTP(chart, displayValueFiles, values, options)
		},
	}

	return cmd.Run(ctx, os.Args)
}
