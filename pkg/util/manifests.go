package util

import (
	"slices"
	"strings"

	"knit/pkg/types"
)

func SortManifests(manifests []types.Manifest) {
	slices.SortFunc(manifests, func(a, b types.Manifest) int {
		return compareManifests(a, b,
			func(x types.Manifest) string { return x["apiVersion"].(string) },
			func(x types.Manifest) string { return x["kind"].(string) },
			func(x types.Manifest) string {
				metadata := x["metadata"].(map[string]any)
				return metadata["name"].(string)
			},
			func(x types.Manifest) string {
				metadata := x["metadata"].(map[string]any)
				return metadata["generateName"].(string)
			})
	})
}

func compareManifests(a, b types.Manifest, selectors ...func(x types.Manifest) string) int {
	var comparator int
	for _, selector := range selectors {
		comparator = strings.Compare(selector(a), selector(b))
		if comparator != 0 {
			return comparator
		}
	}

	return 0
}
