package plugin

import (
	"fmt"
	"knit/pkg/util"

	"gopkg.in/yaml.v3"
	"kcl-lang.io/lib/go/plugin"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func init() {
	plugin.RegisterPlugin(plugin.Plugin{
		Name: "kustomize",
		MethodMap: map[string]plugin.MethodSpec{
			"build": {
				// kustomize.build(kustomization, [resources])
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					kustomizationArg := getCallArgs(args, "kustomization", 0)
					resourceArgs := args.ListArg(1)

					kustomization, ok := kustomizationArg.(map[string]any)
					if !ok {
						return nil, fmt.Errorf("expecting kustomization to be a map with string keys")
					}

					fSys := filesys.MakeFsOnDisk()

					tmpDir, err := util.NewTempDir("kustomize")
					if err != nil {
						return nil, err
					}
					defer tmpDir.Remove()

					err = appendAdditionalResources(kustomization, tmpDir, fSys, resourceArgs)
					if err != nil {
						return nil, err
					}

					kustomizationContent, err := yaml.Marshal(&kustomization)
					if err != nil {
						return nil, err
					}
					kustomizationPath, err := tmpDir.CreatePath("kustomization.yaml")
					if err != nil {
						return nil, err
					}
					if err := fSys.WriteFile(kustomizationPath, []byte(kustomizationContent)); err != nil {
						return nil, fmt.Errorf("error writing kustomization.yaml: %v", err)
					}

					k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
					resMap, err := k.Run(fSys, tmpDir.Path)
					if err != nil {
						return nil, fmt.Errorf("error running kustomize: %v", err)
					}

					return &plugin.MethodResult{V: resMap.Resources()}, nil
				},
			},
		},
	})
}

func appendAdditionalResources(kustomization map[string]any, tmpDir *util.TempDir, fSys filesys.FileSystem, resources []any) error {
	var newResources []string
	oldResources, ok := kustomization["resources"]
	if ok {
		var err error
		newResources, err = util.AnySliceToTyped[string](oldResources)
		if err != nil {
			return fmt.Errorf("resources must be a list of strings: %v", err)
		}
	}

	for i, resource := range resources {
		resourceYAML, err := yaml.Marshal(resource)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("resource-%d.yaml", i)
		resourcePath, err := tmpDir.CreatePath(fileName)
		if err != nil {
			return err
		}
		if err := fSys.WriteFile(resourcePath, []byte(resourceYAML)); err != nil {
			return err
		}

		newResources = append(newResources, fileName)
	}

	kustomization["resources"] = newResources

	return nil
}
