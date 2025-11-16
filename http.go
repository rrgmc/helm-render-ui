package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"path"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"
)

const devHTTPPort = 17821

func runHTTP(ctx context.Context, httpPort int, chart *chart.Chart, valueFiles []string, values map[string]any, releaseOptions chartutil.ReleaseOptions) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/data", httpHandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		fnprefix := fmt.Sprintf("%s/templates/", chart.Name())

		valuesToRender, err := chartutil.ToRenderValues(chart, values, releaseOptions, nil)
		if err != nil {
			return err
		}

		renderedTemplate, err := engine.Render(chart, valuesToRender)
		if err != nil {
			return fmt.Errorf("cannot render template using engine: %v", err)
		}

		chartStr, err := yaml.Marshal(chart.Metadata)
		if err != nil {
			return err
		}

		chartStrValue := string(chartStr)

		if len(valueFiles) > 0 {
			chartStrValue += "\n---\nvalue_files:\n"
			for _, file := range valueFiles {
				chartStrValue += fmt.Sprintf("- %s\n", strings.TrimPrefix(file, fnprefix))
			}
		}

		releaseStr, err := yaml.Marshal(releaseOptions)
		if err != nil {
			return err
		}

		valuesStr, err := yaml.Marshal(values)
		if err != nil {
			return err
		}

		fullValuesStr, err := yaml.Marshal(valuesToRender["Values"])
		if err != nil {
			return err
		}

		renderValuesStr, err := yaml.Marshal(valuesToRender)
		if err != nil {
			return err
		}

		data := apiData{
			Chart:        chartStrValue,
			Release:      string(releaseStr),
			Values:       string(valuesStr),
			FullValues:   string(fullValuesStr),
			RenderValues: string(renderValuesStr),
		}

		for cf := range chartFilesIter(chart) {
			fv, ok := renderedTemplate[cf.FullPath]
			if !ok {
				// slog.Warn("cannot find rendered template", "template", cf.FullPath)
				continue
			}

			if strings.TrimSpace(fv) == "" {
				continue
			}

			fileDesc := path.Join(cf.Path...)
			if len(cf.Path) > 0 {
				fileDesc += "/"
			}
			fileDesc += cf.Filename

			data.PreviewFiles = append(data.PreviewFiles, apiDataFile{
				Filename: fileDesc,
				Preview:  fv,
			})
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		return json.NewEncoder(w).Encode(data)
	}))
	err := uiHandler(mux)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close() // Ensure the listener is closed when main exits

	serverHTTPPort := listener.Addr().(*net.TCPAddr).Port

	browserURL := fmt.Sprintf("http://127.0.0.1:%d", serverHTTPPort)
	if httpPort == 0 {
		slog.InfoContext(ctx, "opening browser URL", "url", browserURL)
		_ = openURL(browserURL)
	} else {
		slog.InfoContext(ctx, "browser URL", "url", browserURL)
	}

	return http.Serve(listener, mux)
}

func httpHandlerWithError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			slog.ErrorContext(r.Context(), "http handler error: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
