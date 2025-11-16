package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rrgmc/helm-render-ui/helm"
	"github.com/urfave/cli/v3"
	"helm.sh/helm/v3/pkg/chart"
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
				Name:  "repo",
				Usage: "helm repository URL. If set, the folder name parameter will be used as the chart name",
			},
			&cli.StringFlag{
				Name:  "chart-version",
				Usage: "chart version (if downloading from repository)",
			},
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
				Name:    "values",
				Aliases: []string{"f"},
				Usage:   "extra configuration values file name",
			},
			&cli.IntFlag{
				Name:    "http-port",
				Aliases: []string{"p"},
				Usage:   "http port",
			},
			&cli.BoolFlag{
				Name:  "is-upgrade",
				Usage: "sets upgrade mode",
				Value: false,
			},
			&cli.BoolFlag{
				Name:   "dev-port",
				Usage:  "dev http port",
				Hidden: true,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			chartRepo := command.String("repo")
			chartFolder := command.StringArgs("helm-chart-folder")[0]

			if strings.TrimSpace(chartFolder) == "" {
				return fmt.Errorf("helm chart folder is required")
			}
			var err error

			var cht *chart.Chart
			var chartVersions []string
			if chartRepo != "" {
				chartName := chartFolder
				slog.InfoContext(ctx, "loading chart from repository",
					"repo", chartRepo,
					"chart", chartFolder,
					"version", command.String("chart-version"))

				repository, err := helm.LoadRepository(chartRepo)
				if err != nil {
					return err
				}

				latestChart, err := repository.GetChart(chartFolder, command.String("chart-version"))
				if err != nil {
					return err
				}

				latestChartFiles, err := latestChart.Download()
				if err != nil {
					return err
				}
				defer latestChartFiles.Close()

				chartFolder = latestChartFiles.ChartPath()

				for entry, err := range repository.ChartVersions(chartName, 20) {
					if err != nil {
						slog.Warn("error listing chart versions", "error", err)
						break
					}
					var date string
					if !entry.Created.IsZero() {
						date = fmt.Sprintf(" [%s]", entry.Created.Format(time.RFC3339))
					}
					chartVersions = append(chartVersions, fmt.Sprintf("%s%s", entry.Version, date))
				}
			}

			cht, err = loader.LoadDir(chartFolder)
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

			if err := chartutil.ProcessDependencies(cht, values); err != nil {
				return err
			}

			options := chartutil.ReleaseOptions{
				Name:      command.String("release"),
				Namespace: command.String("namespace"),
				Revision:  1,
				IsInstall: !command.Bool("is-upgrade"),
				IsUpgrade: command.Bool("is-upgrade"),
			}
			if options.Name == "" {
				options.Name = cht.Metadata.Name
			}

			httpPort := command.Int("http-port")
			if command.Bool("dev-port") {
				httpPort = devHTTPPort
			}

			return runHTTP(ctx, httpPort, cht, displayValueFiles, values, options, chartVersions)
		},
	}

	return cmd.Run(ctx, os.Args)
}
