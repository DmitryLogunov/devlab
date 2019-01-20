package helpers

import (
	"devlab/lib/files"
	"errors"
	"path/filepath"
	"sort"
	"strconv"

	contextErrors "devlab/cmd/context/common/errors"
)

// GetConfiguration reads, parses and returns configuration yaml file
func GetConfiguration(config map[string]map[string]string) (configuration map[string]map[string]string, err error) {
	if config["paths"]["configurations"] == "" || config["description"]["configuration"] == "" {
		return make(map[string]map[string]string), contextErrors.ErrNotDefinedConfigurationPath
	}

	configurationTemplatesPath := config["paths"]["configurations"] + "/" + config["description"]["configuration"]
	if isConfigurationExists, _ := files.IsExists("./" + configurationTemplatesPath); !isConfigurationExists {
		return make(map[string]map[string]string), contextErrors.ErrCouldntReadConfiguration
	}

	configuration, err = files.ReadTwoLevelYaml(configurationTemplatesPath)
	if err != nil {
		return make(map[string]map[string]string), contextErrors.ErrCouldntParseConfiguration
	}

	configurationTemlateDir := filepath.Dir(configurationTemplatesPath)
	defaultConfiguratonPath := configurationTemlateDir + "/default.yml"
	if isDefaultConfigurationExists, _ := files.IsExists(defaultConfiguratonPath); isDefaultConfigurationExists {
		defaultConfiguration, err := files.ReadTwoLevelYaml(defaultConfiguratonPath)
		if err != nil {
			return make(map[string]map[string]string), contextErrors.ErrCouldntParseConfiguration
		}
		configuration = MergeMaps(configuration, defaultConfiguration)
	}

	return
}

// GetValueFromContextOrDefault returns value form context map accordingly topLevel ans subLevel keys.
// If context map value is empty then return value from default map
func GetValueFromContextOrDefault(context map[string]map[string]string,
	defaultConfig map[string]map[string]string, topLevelKey, subLevelKey string) (value string) {

	if context[topLevelKey][subLevelKey] != "" {
		return context[topLevelKey][subLevelKey]
	}
	return defaultConfig[topLevelKey][subLevelKey]
}

// MergeMaps merges source map and default map with priority of source map
func MergeMaps(sourceMap, defaultMap map[string]map[string]string) (mergedMap map[string]map[string]string) {
	mergedMap = defaultMap
	for itemKey, itemMap := range sourceMap {
		for key, value := range itemMap {
			if value != "" {
				mergedMap[itemKey][key] = value
			}
		}
	}

	return
}

// KeyWithOrder implements type key with order
type KeyWithOrder struct {
	Key   string
	order int
}

// GetSortedKeysFromYaml returns sorted keys from yaml
//
// Example fo yaml:
//
// key1:
//   order: 1
// key2:
//   oreder: 0
// key3:
//   oreder: 2
// ....
// It should return : [key2, key1, key3]
func GetSortedKeysFromYaml(path string) (keys []KeyWithOrder, err error) {
	data, err := files.ReadTwoLevelYaml(path)
	if err != nil {
		return
	}

	for key, orderMap := range data {
		order, _ := strconv.Atoi(orderMap["order"])
		if key != "" {
			orderedKey := KeyWithOrder{key, order}
			keys = append(keys, orderedKey)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].order > keys[j].order
	})

	return
}

// GetMainConfig returns main config
func GetMainConfig() (config map[string]map[string]string, err error) {
	config, err = files.ReadMainConfig()
	if err != nil {
		return make(map[string]map[string]string), contextErrors.ErrCouldntReadConfig
	}
	if config["paths"]["contexts"] == "" {
		return make(map[string]map[string]string), contextErrors.ErrNotDefinedContextsPath
	}
	return
}

// GetApplicationServices returns application services settings by merging default template settings
// and context level settings
func GetApplicationServices(contextName string, context, config map[string]map[string]string,
	configuration map[string]map[string]string) (applicationServices map[string]map[string]string, err error) {

	applicationServices = make(map[string]map[string]string)

	templatesPath := GetValueFromContextOrDefault(context, config, "paths", "context-templates")
	applicationServicesTemplate := GetValueFromContextOrDefault(context, configuration, "application-services", "template")
	applicationServicesTemplatePath := "./" + templatesPath + "/application-services/" + applicationServicesTemplate

	isAplicationServicesDirExists, _ := files.IsExists(applicationServicesTemplatePath)
	if !isAplicationServicesDirExists {
		err = errors.New("applications services directory does not exist")
		return
	}

	applicationServicesFromTemplate, err := files.ReadTwoLevelYaml(applicationServicesTemplatePath)
	if err != nil {
		return
	}

	// applications services settings equal template settings as default
	applicationServices = applicationServicesFromTemplate

	// checking if application-services settings from context exist and merge its if yes
	levelFolder := config["configuration-levels"]["applications"]
	applicationServicesContextSettingsPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/" + levelFolder + "/application-services.settings.yml"
	isAplicationServicesContextSettingsExists, _ := files.IsExists(applicationServicesContextSettingsPath)
	if isAplicationServicesContextSettingsExists {
		applicationServicesFromContext, _ := files.ReadTwoLevelYaml(applicationServicesContextSettingsPath)
		applicationServices = MergeMaps(applicationServicesFromContext, applicationServicesFromTemplate)
	}

	return
}
