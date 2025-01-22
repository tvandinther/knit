package plugin

import (
	"fmt"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/releaseutil"
	"kcl-lang.io/lib/go/plugin"

	"knit/pkg/helm"
	"knit/pkg/types"
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
					err := util.MapToStruct(chartArgMap, &chart)
					if err != nil {
						return nil, fmt.Errorf("unable to map chart to struct: %w", err)
					}

					err = validate(&chart)
					if err != nil {
						return nil, fmt.Errorf("invalid chart: %w", err)
					}

					release, err := helm.RunTemplate(&helm.ChartRef{
						Repository: chart.Repository,
						Name:       chart.Name,
						Version:    chart.Version,
					}, chart.Values, chart.ReleaseName, chart.Namespace, chart.Capabilities.APIVersions)
					if err != nil {
						return nil, fmt.Errorf("problem templating helm chart: %w", err)
					}

					splitManifests := releaseutil.SplitManifests(release.Manifest)
					var manifestSlice []types.Manifest
					for _, manifestString := range splitManifests {
						var manifest types.Manifest
						yaml.Unmarshal([]byte(manifestString), &manifest)
						if len(manifest) > 0 {
							manifestSlice = append(manifestSlice, manifest)
						}
					}

					for _, hook := range release.Hooks {
						var manifest types.Manifest
						yaml.Unmarshal([]byte(hook.Manifest), &manifest)
						if len(manifest) > 0 {
							manifestSlice = append(manifestSlice, manifest)
						}
					}

					util.SortManifests(manifestSlice)

					return &plugin.MethodResult{V: manifestSlice}, nil
				},
			},
		},
	})
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
