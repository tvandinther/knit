package helm

import (
	"bytes"
	"fmt"
	"knit/pkg/logging"
	"log"
	"sort"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
)

func getChart(chartRef *ChartRef, settings *cli.EnvSettings) (*chart.Chart, error) {
	logger := logging.GetInstance()

	cpo := action.ChartPathOptions{
		RepoURL: chartRef.Repository,
		Version: chartRef.Version,
	}

	var chartPath string
	var err error
	locateChart := func() {
		chartPath, err = cpo.LocateChart(chartRef.Name, settings)
	}

	locateChart()
	if err != nil {
		err := pullChart(logger, settings, chartRef)
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to pull chart: ", err))
			return nil, err
		}	

		locateChart()
		if err != nil {
			log.Fatalln(fmt.Errorf("failed to locate chart: ", err))
			return nil, err
		}
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to load chart: ", err))
		return nil, err
	}

	return chart, nil
}

func RunTemplate(chartRef *ChartRef) error {
	settings := cli.New()
	chart, err := getChart(chartRef, settings)

	registryClient, err := newRegistryClient(settings, false)
	if err != nil {
		return fmt.Errorf("failed to created registry client: %w", err)
	}

	kubeVersion, err := chartutil.ParseKubeVersion("1.31.0")
	if err != nil {
		return err
	}
	
	client := action.NewInstall(&action.Configuration{})
	client.SetRegistryClient(registryClient)
	client.DryRun = true
	client.ClientOnly = true
	client.ReleaseName = "todo" // TODO
	client.KubeVersion = kubeVersion
	client.APIVersions = chartutil.DefaultVersionSet // TODO: What's this?
	client.IncludeCRDs = true

	rel, err := client.Run(chart, nil)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to run install: ", err))
		return err
	}

	if rel == nil {
		return fmt.Errorf("nil release")
	}

	printRelease(rel)
	
	return nil
}

func printRelease(rel *release.Release) {
	if rel != nil {
		var manifests bytes.Buffer
		fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))
		// if !client.DisableHooks {
		// 	fileWritten := make(map[string]bool)
		// 	for _, m := range rel.Hooks {
		// 		if skipTests && isTestHook(m) {
		// 			continue
		// 		}
		// 		if client.OutputDir == "" {
		// 			fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
		// 		} else {
		// 			newDir := client.OutputDir
		// 			if client.UseReleaseName {
		// 				newDir = filepath.Join(client.OutputDir, client.ReleaseName)
		// 			}
		// 			_, err := os.Stat(filepath.Join(newDir, m.Path))
		// 			if err == nil {
		// 				fileWritten[m.Path] = true
		// 			}

		// 			err = writeToFile(newDir, m.Path, m.Manifest, fileWritten[m.Path])
		// 			if err != nil {
		// 				return err
		// 			}
		// 		}

		// 	}
		// }

		// if we have a list of files to render, then check that each of the
		// provided files exists in the chart.
		// if len(showFiles) > 0 {
		if false {
			// This is necessary to ensure consistent manifest ordering when using --show-only
			// with globs or directory names.
			splitManifests := releaseutil.SplitManifests(manifests.String())
			manifestsKeys := make([]string, 0, len(splitManifests))
			for k := range splitManifests {
				manifestsKeys = append(manifestsKeys, k)
			}
			sort.Sort(releaseutil.BySplitManifestsOrder(manifestsKeys))

			// manifestNameRegex := regexp.MustCompile("# Source: [^/]+/(.+)")
			var manifestsToRender []string
			// for _, f := range showFiles {
			// 	missing := true
			// 	// Use linux-style filepath separators to unify user's input path
			// 	f = filepath.ToSlash(f)
			// 	for _, manifestKey := range manifestsKeys {
			// 		manifest := splitManifests[manifestKey]
			// 		submatch := manifestNameRegex.FindStringSubmatch(manifest)
			// 		if len(submatch) == 0 {
			// 			continue
			// 		}
			// 		manifestName := submatch[1]
			// 		// manifest.Name is rendered using linux-style filepath separators on Windows as
			// 		// well as macOS/linux.
			// 		manifestPathSplit := strings.Split(manifestName, "/")
			// 		// manifest.Path is connected using linux-style filepath separators on Windows as
			// 		// well as macOS/linux
			// 		manifestPath := strings.Join(manifestPathSplit, "/")

			// 		// if the filepath provided matches a manifest path in the
			// 		// chart, render that manifest
			// 		if matched, _ := filepath.Match(f, manifestPath); !matched {
			// 			continue
			// 		}
			// 		manifestsToRender = append(manifestsToRender, manifest)
			// 		missing = false
			// 	}
			// 	if missing {
			// 		return fmt.Errorf("could not find template %s in chart", f)
			// 	}
			// }
			for _, m := range manifestsToRender {
				out := log.Writer()
				fmt.Fprintf(out, "---\n%s\n", m)
			}
		} else {
			out := log.Writer()
			fmt.Fprintf(out, "%s", manifests.String())
		}
	}
}
