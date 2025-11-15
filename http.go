package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"
)

const httpPort = "17821"

func runHTTP(ctx context.Context, chart *chart.Chart, values map[string]any, releaseOptions chartutil.ReleaseOptions) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/data", httpHandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
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
			Chart:        string(chartStr),
			Release:      string(releaseStr),
			Values:       string(valuesStr),
			RenderValues: string(renderValuesStr),
			Preview:      outputTemplate(renderedTemplate),
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		return json.NewEncoder(w).Encode(data)
	}))

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
