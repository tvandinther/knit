package plugin

import (
	"fmt"

	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/releaseutil"
	"kcl-lang.io/lib/go/plugin"

	"knit/pkg/helm"
	"knit/pkg/util"
)

type ChartArg struct {
	Repository string
	Name       string
	Version    string
	Values     map[string]interface{}
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
					util.MapToStruct(chartArgMap, &chart)

					release, err := helm.RunTemplate(&helm.ChartRef{
						Repository: chart.Repository,
						Name:       chart.Name,
						Version:    chart.Version,
					}, chart.Values)
					if err != nil {
						return nil, err
					}

					splitManifests := releaseutil.SplitManifests(release.Manifest)

					var manifestSlice []interface{}
					for _, manifestString := range splitManifests {
						var manifest interface{}
						yaml.Unmarshal([]byte(manifestString), &manifest)
						manifestSlice = append(manifestSlice, manifest)
					}

					return &plugin.MethodResult{V: manifestSlice}, nil
				},
			},
		},
	})
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
