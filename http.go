package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

const httpPort = "17821"

func runHTTP(ctx context.Context, chart *chart.Chart, values map[string]any, releaseOptions chartutil.ReleaseOptions) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		valuesToRender, err := chartutil.ToRenderValues(chart, values, releaseOptions, nil)
		if err != nil {
			return err
		}

		renderedTemplate, err := engine.Render(chart, valuesToRender)
		if err != nil {
			return fmt.Errorf("cannot render template using engine: %v", err)
		}

		_, err = w.Write([]byte(outputTemplate(renderedTemplate)))
		return err
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
