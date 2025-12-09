package main

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"os/exec"
	"path"
	"runtime"
	"slices"
	"strings"

	"helm.sh/helm/v3/pkg/chart"
)

func ensureRelativePath(path string) string {
	return strings.TrimLeft(path, "/")
}

func formatHelmFilename(ch *chart.Chart, name string) string {
	name = strings.TrimPrefix(name, ch.Name())
	name = strings.TrimPrefix(name, "/templates")
	return ensureRelativePath(name)
}

func outputTemplate(renderedTemplate map[string]string) string {
	var tmpl strings.Builder
	for fn, fv := range mapSortedByKey(renderedTemplate) {
		if strings.TrimSpace(fv) == "" {
			continue
		}
		tmpl.WriteString(fmt.Sprintf("# %s\n", fn))
		tmpl.WriteString(ensureNewline(fv))
		tmpl.WriteString("---\n")
	}

	return tmpl.String()
}

// mapSortedByKey returns an iterator for the given map that
// yields the key-value pairs in sorted order.
func mapSortedByKey[Map ~map[K]V, K cmp.Ordered, V any](m Map) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.Sorted(maps.Keys(m)) {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

type chartIterData struct {
	FullPath string
	Path     []string
	Filename string
}

func chartFilesIter(ch *chart.Chart) iter.Seq[chartIterData] {
	return func(yield func(chartIterData) bool) {
		for _, tmpl := range ch.Templates {
			cid := chartIterData{
				FullPath: path.Join(ch.ChartFullPath(), tmpl.Name),
				Filename: strings.TrimPrefix(tmpl.Name, "templates/"),
			}
			cp := ch
			for cp != nil {
				if cp.Parent() == nil {
					break
				}
				cid.Path = append(cid.Path, cp.Name())
				cp = cp.Parent()
			}
			slices.Reverse(cid.Path)

			if !yield(cid) {
				return
			}
		}
		for _, dep := range ch.Dependencies() {
			for ct := range chartFilesIter(dep) {
				if !yield(ct) {
					return
				}
			}
		}
	}
}

func ensureNewline(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}

// openURL opens the specified URL in the default browser of the user.
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}
	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}
	return exec.Command(cmd, args...).Start()
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
