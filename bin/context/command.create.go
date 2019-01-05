package context

import (
	"devlab/lib/files"
	"devlab/lib/logger"
	"errors"
	"fmt"
)

// Create a new context
func Create(contextName string) (err error) {
	if contextName == "" {
		return ErrNotDefinedContextName
	}

	config, err := files.ReadMainConfig()
	if err != nil {
		return ErrCouldntReadConfig
	}
	if config["paths"]["contexts"] == "" {
		return ErrNotDefinedContextsPath
	}

	// check context dir and create it if need
	contextDir := "./" + config["paths"]["contexts"] + "/" + contextName
	if isContextDirExists, _ := files.IsExists(contextDir); !isContextDirExists {
		files.CreateDir("./" + contextDir)
	}

	contextSettingsParameters, err := getSettingsParamters(contextName, config)
	if err != nil {
		return err
	}

	if config["paths"]["context-templates"] == "" {
		return ErrNotDefinedContextsPath
	}
	defaultSettingsPath := "./" + config["paths"]["context-templates"] + "/default-settings.yml"

	contextSettingsPath := contextDir + "/context.settings.yml"
	isContextSettingsExists, _ := files.IsExists(contextSettingsPath)
	if isContextSettingsExists {
		logger.Warn("There is settings.yml file in context: %s. This file will be rewrited!", contextSettingsPath)
		files.Delete(contextSettingsPath)
	}

	// create context.settings.yml
	err = files.RenderTextTemplate(defaultSettingsPath, contextSettingsPath, contextSettingsParameters)
	if err != nil {
		return err
	}

	// create builing.settings.yml
	err = createContextBlockSettingsFile(contextName, config, "building")
	if err != nil {
		return err
	}

	// create deploying.settings.yml
	err = createContextBlockSettingsFile(contextName, config, "deploying")
	if err != nil {
		return err
	}

	// create system-services.settings.yml
	err = createContextBlockSettingsFile(contextName, config, "system-services")
	if err != nil {
		return err
	}

	// create application-services.settings.yml
	paramsFilter := func(defaultValuesMap map[string]string) func(param string, value *string) {
		return func(param string, value *string) {
			for key, defaultValue := range defaultValuesMap {
				if param == key {
					if *value == "" {
						*value = defaultValue
					}
				}
			}
		}
	}
	configuration, err := getConfiguration(config)
	if err != nil {
		return err
	}
	defaultValues := make(map[string]string)
	defaultValues["branch"] = configuration["application-services"]["base-branch"]
	if configuration["application-services"]["feature-branch-naming"] == "context-name" {
		defaultValues["branch"] = contextName
	}
	for param, value := range configuration["application-services"] {
		if param != "base-branch" && param != "feature-branch" && param != "template" {
			defaultValues[param] = value
		}
	}
	err = createContextBlockSettingsFile(contextName, config, "application-services", paramsFilter(defaultValues))
	if err != nil {
		return err
	}

	return
}

// Get settings parameters of context
func getSettingsParamters(contextName string, config map[string]map[string]string) (settingsParameters SettingsParameters, err error) {
	configuration, err := getConfiguration(config)
	if err != nil {
		return SettingsParameters{}, err
	}

	settingsParameters = SettingsParameters{
		Name: contextName,
		Description: SettingsParametersDescription{
			Maintainer: fmt.Sprintf("\"%s\"", config["description"]["maintainer"]),
			Created:    fmt.Sprintf("\"%s\"", "2019-01-01")},
		Docker: SettingsParametersDocker{
			ImagesPrefix: fmt.Sprintf("%s-%s", configuration["docker"]["images-prefix"], contextName),
			Network:      fmt.Sprintf("%s-%s", configuration["docker"]["network-prefix"], contextName)},
		Paths: SettingsParametersPaths{
			Templates: config["paths"]["context-templates"]},
		Building: SettingsParametersBuilding{
			Template: configuration["building"]["template"]},
		Deploying: SettingsParametersDeploying{
			Template: configuration["deploying"]["template"]},
		SystemServices: SettingsParametersSystemServices{
			Template: configuration["system-services"]["template"]},
		ApplicationServices: SettingsParametersApplicationServices{
			BaseBranch:          configuration["application-services"]["base-branch"],
			FeatureBranchNaming: configuration["application-services"]["feature-branch-naming"],
			Template:            configuration["application-services"]["template"]}}

	return
}

// Create block of context settings yaml file
type filterParamsFunc func(string, *string)

func createContextBlockSettingsFile(contextName string, config map[string]map[string]string,
	blockName string, filterParams ...filterParamsFunc) (err error) {

	configuration, err := getConfiguration(config)
	if err != nil {
		return err
	}

	if configuration[blockName]["template"] == "" {
		errInfo := fmt.Sprintf("ERROR: context => Template for block '%s' is not defined. Couldn't create file settings for this block", blockName)
		return errors.New(errInfo)
	}

	// read context settings
	contextSettingsPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/context.settings.yml"
	context, err := files.ReadTwoLevelYaml(contextSettingsPath)
	if err != nil {
		return err
	}

	// set template path
	templatesPath := getValueFromContextOrDefault(context, config, "paths", "context-templates")

	// set context block settings template
	contextBlockTemplate := getValueFromContextOrDefault(context, configuration, blockName, "template")
	contextBlockTemplatePath := "./" + templatesPath + "/" + blockName + "/" + contextBlockTemplate
	contextBlockTemplateSettings, err := files.ReadTwoLevelYaml(contextBlockTemplatePath)
	if err != nil {
		return err
	}

	// read parameters of context block
	parametersPath := "./" + templatesPath + "/" + blockName + "/parameters.yml"
	parametersOfBlockItems, err := getSortedKeysFromYaml(parametersPath)
	if err != nil {
		return err
	}

	// check if context block settings file exists and delete it if yes
	contextBlockSettingsFilePath := "./" + config["paths"]["contexts"] + "/" + contextName + "/" + blockName + ".settings.yml"
	isCcontextBlockSettingsFileExists, _ := files.IsExists(contextBlockSettingsFilePath)
	if isCcontextBlockSettingsFileExists {
		files.Delete(contextBlockSettingsFilePath)
	}

	// fill context block settings map
	for item := range contextBlockTemplateSettings {
		_, err = files.WriteAppendFileWithIndent(contextBlockSettingsFilePath, item+": ", 0)
		if err != nil {
			return err
		}
		for i := 0; i < len(parametersOfBlockItems); i++ {
			param := parametersOfBlockItems[i].key
			value := ""
			if contextBlockTemplateSettings[item][param] != "" {
				value = contextBlockTemplateSettings[item][param]
			}
			if filterParams != nil && filterParams[0] != nil {
				filterParams[0](param, &value)
			}

			_, err = files.WriteAppendFileWithIndent(contextBlockSettingsFilePath, param+": "+value, 2)
			if err != nil {
				return err
			}
		}
	}

	return
}
