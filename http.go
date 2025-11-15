package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"
)

const httpPort = "17821"

func runHTTP(chart *chart.Chart, valueFiles []string, values map[string]any, releaseOptions chartutil.ReleaseOptions) error {
	uiFS, err := fs.Sub(staticFS, "ui/build")
	if err != nil {
		return err
	}

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

		chartStrValue := string(chartStr) + "\n---\nvalue_files:\n"
		for _, file := range valueFiles {
			chartStrValue += fmt.Sprintf("- %s\n", strings.TrimPrefix(file, fnprefix))
		}

		releaseStr, err := yaml.Marshal(releaseOptions)
		if err != nil {
			return err
		}

		valuesStr, err := yaml.Marshal(values)
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
			RenderValues: string(renderValuesStr),
			// Preview:      outputTemplate(renderedTemplate),
		}

		for fn, fv := range mapSortedByKey(renderedTemplate) {
			if strings.TrimSpace(fv) == "" {
				continue
			}
			data.PreviewFiles = append(data.PreviewFiles, apiDataFile{
				Filename: strings.TrimPrefix(fn, fnprefix),
				Preview:  fv,
			})
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		return json.NewEncoder(w).Encode(data)
	}))
	mux.Handle("/", http.FileServer(http.FS(uiFS)))

	return http.ListenAndServe(":"+httpPort, mux)
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
