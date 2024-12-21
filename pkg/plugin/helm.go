package plugin

import (
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/releaseutil"
	"kcl-lang.io/lib/go/plugin"

	"knit/pkg/helm"
	"knit/pkg/util"
)

type ChartArg struct {
	Repository   string
	Name         string
	Version      string
	ReleaseName  string
	Namespace    string
	Values       map[string]interface{}
	Capabilities Capabilities
}

type Capabilities struct {
	APIVersions []string
}

type Manifest = map[string]any

func init() {
	plugin.RegisterPlugin(plugin.Plugin{
		Name: "helm",
		MethodMap: map[string]plugin.MethodSpec{
			"template": {
				// helm.template(chart)
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					chartArg := getCallArgs(args, "chart", 0)

					chartArgMap, ok := chartArg.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("invalid argument: %+v", chartArg)
					}

					var chart ChartArg
					util.MapToStruct(chartArgMap, &chart)

					err := validate(&chart)
					if err != nil {
						return nil, err
					}

					release, err := helm.RunTemplate(&helm.ChartRef{
						Repository: chart.Repository,
						Name:       chart.Name,
						Version:    chart.Version,
					}, chart.Values, chart.ReleaseName, chart.Namespace, chart.Capabilities.APIVersions)
					if err != nil {
						return nil, err
					}

					splitManifests := releaseutil.SplitManifests(release.Manifest)

					var manifestSlice []Manifest
					for _, manifestString := range splitManifests {
						var manifest Manifest
						yaml.Unmarshal([]byte(manifestString), &manifest)
						manifestSlice = append(manifestSlice, manifest)
					}

					slices.SortFunc(manifestSlice, func(a, b Manifest) int {
						return sortManifests(a, b,
							func(x Manifest) string { return x["apiVersion"].(string) },
							func(x Manifest) string { return x["kind"].(string) },
							func(x Manifest) string {
								metadata := x["metadata"].(map[string]any)
								return metadata["name"].(string)
							},
							func(x Manifest) string {
								metadata := x["metadata"].(map[string]any)
								return metadata["generateName"].(string)
							})
					})

					return &plugin.MethodResult{V: manifestSlice}, nil
				},
			},
		},
	})
}

func sortManifests(a, b Manifest, selectors ...func(x Manifest) string) int {
	var comparator int
	for _, selector := range selectors {
		comparator = strings.Compare(selector(a), selector(b))
		if comparator != 0 {
			return comparator
		}
	}

	return 0
}

func validate(chartArg *ChartArg) error {
	if len(chartArg.ReleaseName) < 1 {
		return fmt.Errorf("releaseName must be defined on the chart")
	}

	return nil
}

func getCallArgs(p *plugin.MethodArgs, key string, index int) any {
	if val, ok := p.KwArgs[key]; ok {
		return val
	}
	if index < len(p.Args) {
		return p.Args[index]
	}
	return nil
}
