package util

import (
	"fmt"
	"slices"
	"strings"

	"knit/pkg/types"
)

func SortManifests(manifests []types.Manifest) error {
	errors := make([]error, 0)
	slices.SortFunc(manifests, func(a, b types.Manifest) int {
		comp, err := compareManifests(a, b,
			func(x types.Manifest) (string, error) { return getFromMap[string](x, "apiVersion") },
			func(x types.Manifest) (string, error) { return getFromMap[string](x, "kind") },
			func(x types.Manifest) (string, error) {
				metadata, err := getFromMap[map[string]any](x, "metadata")
				if err != nil {
					return "", err
				}
				return getFromMap[string](metadata, "name")
			},
			func(x types.Manifest) (string, error) {
				metadata, err := getFromMap[map[string]any](x, "metadata")
				if err != nil {
					return "", err
				}
				return getFromMap[string](metadata, "generateName")
			})
		if err != nil {
			errors = append(errors, err)
		}

		return comp
	})

	if len(errors) != 0 {
		var errorString string
		for _, err := range errors {
			errorString += err.Error()
		}
		return fmt.Errorf("encountered %d errors during sorting\n%s", len(errors), errorString)
	}

	return nil
}

func getFromMap[T any](manifest map[string]any, key string) (T, error) {
	str, ok := manifest[key].(T)
	if !ok {
		var none T
		return none, fmt.Errorf("%s of manifest can not be converted to string\n%+v", key, manifest)
	}

	return str, nil
}

func compareManifests(a, b types.Manifest, selectors ...func(x types.Manifest) (string, error)) (int, error) {
	var comparator int
	for _, selector := range selectors {
		strA, err := selector(a)
		if err != nil {
			return 0, err
		}
		strB, err := selector(b)
		if err != nil {
			return 0, err
		}
		comparator = strings.Compare(strA, strB)
		if comparator != 0 {
			return comparator, nil
		}
	}

	return 0, nil
}
