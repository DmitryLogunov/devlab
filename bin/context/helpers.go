package Context

import (
  "devlab/lib/files"
)

/**
*
*/
func getConfiguration(config map[string]map[string]string) (configuration map[string]map[string]string, err error) {
	if config["paths"]["configurations"] == "" || config["description"]["configuration"] == "" {
    return make(map[string]map[string]string), ErrNotDefinedConfigurationPath  
  }

  configurationTemplatesPath := config["paths"]["configurations"] + "/" + config["description"]["configuration"]
  if isConfigurationExists, _ :=  files.IsExists("./" + configurationTemplatesPath); !isConfigurationExists  { 
    return make(map[string]map[string]string), ErrCouldntReadConfiguration  
  }
  
  configuration, err = files.ReadTwoLevelYaml(configurationTemplatesPath)
  if err != nil {
    return make(map[string]map[string]string), ErrCouldntParseConfiguration
	}
	
	return
}


/**
*
*/
func getValueFromContextOrDefault(context map[string]map[string]string, defaultConfig map[string]map[string]string,
  topLevelKey, subLevelKey string) (value string) {    

  if context[topLevelKey][subLevelKey] != "" { return context[topLevelKey][subLevelKey] }
  return defaultConfig[topLevelKey][subLevelKey]
}