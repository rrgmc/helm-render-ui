package main

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

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

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func ensureNewline(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}
