package main

import (
	"fmt"
	"github.com/saitho/diff-docker-compose/lib"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func exitInvalidArgument(err error) {
	fmt.Println(fmt.Errorf("Cannot continue without loading both files: " + err.Error()))
	fmt.Println("Usage: " + os.Args[0] + " [template file] [actual file]")
	fmt.Println("  [template file] defaults to `docker-compose.yml.template`")
	fmt.Println("    [actual file] defaults to `docker-compose.yml`")
	panic(nil)
}

func main() {
	composeTemplateFilename := "docker-compose.yml.template"
	if len(os.Args) >= 2 && os.Args[1] != "" {
		composeTemplateFilename = os.Args[1]
	}
	composeActualFilename := "docker-compose.yml"
	if len(os.Args) >= 3 && os.Args[2] != "" {
		composeActualFilename = os.Args[2]
	}

	workingDir, _ := os.Getwd()

	if !path.IsAbs(composeTemplateFilename) {
		composeTemplateFilename = path.Join(workingDir, composeTemplateFilename)
	}

	if !path.IsAbs(composeActualFilename) {
		composeActualFilename = path.Join(workingDir, composeActualFilename)
	}

	composeTemplateFile, err := os.ReadFile(composeTemplateFilename)
	if err != nil {
		exitInvalidArgument(err)
	}
	composeActualFile, err := os.ReadFile(composeActualFilename)
	if err != nil {
		exitInvalidArgument(err)
	}

	var composeTemplate map[string]interface{}
	if err := yaml.Unmarshal(composeTemplateFile, &composeTemplate); err != nil {
		exitInvalidArgument(err)
	}
	var composeActual map[string]interface{}
	if err := yaml.Unmarshal(composeActualFile, &composeActual); err != nil {
		exitInvalidArgument(err)
	}

	result := lib.DiffYaml(composeTemplate, composeActual)
	evaluateResults(result)
}

func evaluateResults(result lib.YamlDiffResult) {
	if !result.HasChanged([]string{"services"}) {
		fmt.Println("Services have not changed.")
	} else {
		services := result.Get([]string{"services"})
		var addedServices []string
		var removedServices []string
		var modifiedServices []string
		for _, service := range services {
			serviceName := service.Path[1] // first index is "services", second is service name
			switch service.GetType() {
			case lib.Added:
				addedServices = append(addedServices, serviceName)
				break
			case lib.Removed:
				removedServices = append(removedServices, serviceName)
				break
			case lib.Changed:
				modifiedServices = append(modifiedServices, serviceName)
				break
			}
		}
		fmt.Println("Services locally removed/disabled:")
		for _, service := range removedServices {
			fmt.Println("* " + service)
		}
		fmt.Println("")
		fmt.Println("Services locally added/enabled:")
		for _, service := range addedServices {
			fmt.Println("* " + service)
		}
		fmt.Println("")
		fmt.Println("Services locally modified:")
		for _, service := range modifiedServices {
			fmt.Println("* " + service)
		}
		fmt.Println("")
	}
}
