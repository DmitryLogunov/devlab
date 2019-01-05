package context

import (
	"devlab/lib/files"
	"sort"
	"strconv"
)

// Read, parse and return configuration yaml file
func getConfiguration(config map[string]map[string]string) (configuration map[string]map[string]string, err error) {
	if config["paths"]["configurations"] == "" || config["description"]["configuration"] == "" {
		return make(map[string]map[string]string), ErrNotDefinedConfigurationPath
	}

	configurationTemplatesPath := config["paths"]["configurations"] + "/" + config["description"]["configuration"]
	if isConfigurationExists, _ := files.IsExists("./" + configurationTemplatesPath); !isConfigurationExists {
		return make(map[string]map[string]string), ErrCouldntReadConfiguration
	}

	configuration, err = files.ReadTwoLevelYaml(configurationTemplatesPath)
	if err != nil {
		return make(map[string]map[string]string), ErrCouldntParseConfiguration
	}

	return
}

//  Getting value form context map. If context map value is empty then return value from default map
func getValueFromContextOrDefault(context map[string]map[string]string, defaultConfig map[string]map[string]string,
	topLevelKey, subLevelKey string) (value string) {

	if context[topLevelKey][subLevelKey] != "" {
		return context[topLevelKey][subLevelKey]
	}
	return defaultConfig[topLevelKey][subLevelKey]
}

//  Merge source map and default map with priority of source map
func mergeMaps(sourceMap, defaultMap map[string]map[string]string) (mergedMap map[string]map[string]string) {
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

// Get sorted keys from yaml
// key1:
//   order: 1
// key2:
//   oreder: 0
// key3:
//   oreder: 2
// ....
// It should return : [key2, key1, key3]
// ....
type KeyWithOrder struct {
	key   string
	order int
}

func getSortedKeysFromYaml(path string) (keys []KeyWithOrder, err error) {
	data, err := files.ReadTwoLevelYaml(path)
	if err != nil {
		return
	}

	for key, orderMap := range data {
		order, _ := strconv.Atoi(orderMap["order"])
		if key != "" {
			keyWithOrder := KeyWithOrder{key, order}
			keys = append(keys, keyWithOrder)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].order > keys[j].order
	})

	return
}
