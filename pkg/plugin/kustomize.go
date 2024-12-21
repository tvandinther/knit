package plugin

import (
	"fmt"
	"log"

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
				// kustomize.build(kustomization, ...resources)
				Body: func(args *plugin.MethodArgs) (*plugin.MethodResult, error) {
					kustomizationArg := getCallArgs(args, "kustomization", 0)
					resourceArgs := args.ListArg(1)

					kustomizationYAML, err := yaml.Marshal(kustomizationArg)
					if err != nil {
						log.Fatalf("Error marshalling kustomization: %v", err)
					}

					resourceYAMLs := make([]string, len(resourceArgs))
					for i, resource := range resourceArgs {
						resourceYAML, err := yaml.Marshal(resource)
						if err != nil {
							log.Fatalf("Error marshalling resource %d: %v", i, err)
						}
						resourceYAMLs[i] = string(resourceYAML)
					}

					fSys := filesys.MakeFsInMemory()

					for i, resourceYAML := range resourceYAMLs {
						fileName := fmt.Sprintf("resource-%d.yaml", i)
						if err := fSys.WriteFile(fileName, []byte(resourceYAML)); err != nil {
							log.Fatalf("Error writing resource %d to filesystem: %v", i, err)
						}
					}

					kustomizationContent := "resources:\n"
					for i := range resourceYAMLs {
						kustomizationContent += fmt.Sprintf("- resource-%d.yaml\n", i)
					}
					kustomizationContent += string(kustomizationYAML)
					if err := fSys.WriteFile("kustomization.yaml", []byte(kustomizationContent)); err != nil {
						log.Fatalf("Error writing kustomization.yaml: %v", err)
					}

					k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
					resMap, err := k.Run(fSys, ".")
					if err != nil {
						log.Fatalf("Error running kustomize: %v", err)
					}

					return &plugin.MethodResult{V: resMap.Resources()}, nil
				},
			},
		},
	})
}
