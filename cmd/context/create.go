package context

import (
	"devlab/lib/files"
	"devlab/lib/logger"
	"errors"
	"fmt"

	contextErrors "devlab/cmd/context/common/errors"
	contextHelpers "devlab/cmd/context/common/helpers"
	contextTypes "devlab/cmd/context/common/types"
)

// Create creates a new context
func Create(contextName string) (err error) {
	config, contextSettingsPath, defaultSettingsPath,
		contextSettingsParameters, err := initContextToCreate(contextName)
	if err != nil {
		return
	}

	isContextSettingsExists, _ := files.IsExists(contextSettingsPath)
	if isContextSettingsExists {
		logger.Warn("There is settings.yml file in context: %s. This file will be rewrited!", contextSettingsPath)
		files.Delete(contextSettingsPath)
	}

	// creates context.settings.yml
	err = files.RenderTextTemplate(defaultSettingsPath, contextSettingsPath, contextSettingsParameters)
	if err != nil {
		return
	}

	// creates builing.settings.yml
	err = createContextBlockSettingsFile(contextName, config, "building")
	if err != nil {
		return
	}

	// creates deploying.settings.yml
	err = createContextBlockSettingsFile(contextName, config, "deploying")
	if err != nil {
		return
	}

	// creates system-services.settings.yml
	err = createContextBlockSettingsFile(contextName, config, "system-services")
	if err != nil {
		return
	}

	// creates application-services.settings.yml
	err = createApllicationsServicesContextSettingsFile(contextName, config)

	return
}

// @private

// initContextToCreat initializes and returns main context parameters for Context creation
func initContextToCreate(contextName string) (config map[string]map[string]string,
	contextSettingsPath, defaultSettingsPath string,
	contextSettingsParameters contextTypes.SettingsParameters,
	err error) {

	if contextName == "" {
		err = contextErrors.ErrNotDefinedContextName
		return
	}

	config, err = contextHelpers.GetMainConfig()
	if err != nil {
		return
	}

	contextDir := "./" + config["paths"]["contexts"] + "/" + contextName
	if isContextDirExists, _ := files.IsExists(contextDir); !isContextDirExists {
		files.CreateDir("./" + contextDir)
	}

	contextSettingsPath = contextDir + "/context.settings.yml"
	defaultSettingsPath = "./" + config["paths"]["context-templates"] + "/default-settings.yml"

	contextSettingsParameters, err = getSettingsParamters(contextName, config)
	if err != nil {
		return
	}

	if config["paths"]["context-templates"] == "" {
		err = contextErrors.ErrNotDefinedContextsPath
		return
	}

	return
}

// filterParamsFunc implemets handler type for filtering settings parameters
type filterParamsFunc func(string, *string)

// createContextBlockSettingsFile creates '<blockName>.settings.yml' file with settings of context block
func createContextBlockSettingsFile(contextName string, config map[string]map[string]string,
	blockName string, filterParams ...filterParamsFunc) (err error) {

	configuration, err := contextHelpers.GetConfiguration(config)
	if err != nil {
		return
	}

	if configuration[blockName]["template"] == "" {
		errInfo := fmt.Sprintf("context: template for block '%s' is not defined; couldn't create file settings for this block", blockName)
		return errors.New(errInfo)
	}

	contextSettingsPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/context.settings.yml"
	context, err := files.ReadTwoLevelYaml(contextSettingsPath)
	if err != nil {
		return
	}

	templatesPath := contextHelpers.GetValueFromContextOrDefault(context, config, "paths", "context-templates")

	contextBlockTemplate := contextHelpers.GetValueFromContextOrDefault(context, configuration, blockName, "template")
	contextBlockTemplatePath := "./" + templatesPath + "/" + blockName + "/" + contextBlockTemplate
	contextBlockTemplateSettings, err := files.ReadTwoLevelYaml(contextBlockTemplatePath)
	if err != nil {
		return err
	}

	parametersPath := "./" + templatesPath + "/" + blockName + "/parameters.yml"
	parametersOfBlockItems, err := contextHelpers.GetSortedKeysFromYaml(parametersPath)
	if err != nil {
		return err
	}

	contextBlockSettingsFilePath := "./" + config["paths"]["contexts"] + "/" + contextName + "/" + blockName + ".settings.yml"
	isCcontextBlockSettingsFileExists, _ := files.IsExists(contextBlockSettingsFilePath)
	if isCcontextBlockSettingsFileExists {
		files.Delete(contextBlockSettingsFilePath)
	}

	// fill context block settings map
	for item := range contextBlockTemplateSettings {
		_, err = files.WriteAppendFileWithIndent(contextBlockSettingsFilePath, item+": ", 0)
		if err != nil {
			return
		}

		for i := 0; i < len(parametersOfBlockItems); i++ {
			param := parametersOfBlockItems[i].Key
			value := ""

			if contextBlockTemplateSettings[item][param] != "" {
				value = contextBlockTemplateSettings[item][param]
			}

			if filterParams != nil && filterParams[0] != nil {
				filterParams[0](param, &value)
			}

			_, err = files.WriteAppendFileWithIndent(contextBlockSettingsFilePath, param+": "+value, 2)
			if err != nil {
				return
			}
		}
	}

	return
}

// paramsFilter() returns the function which filters parameters of settings map and
// changes it accordance with defaultValuesMap
func paramsFilter(defaultValuesMap map[string]string) func(param string, value *string) {
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

// createContextBlockSettingsFile creates 'application-services.settings.yml' file with settings of context block
func createApllicationsServicesContextSettingsFile(contextName string,
	config map[string]map[string]string) (err error) {

	configuration, err := contextHelpers.GetConfiguration(config)
	if err != nil {
		return
	}

	defaultValues := make(map[string]string)
	defaultValues["branch"] = configuration["application-services"]["base-branch"]
	if configuration["application-services"]["feature-branch-naming"] == "context-name" {
		defaultValues["branch"] = contextName
	}

	for param, value := range configuration["application-services"] {
		if param != "base-branch" && param != "feature-branch-naming" && param != "template" {
			defaultValues[param] = value
		}
	}

	err = createContextBlockSettingsFile(contextName, config, "application-services", paramsFilter(defaultValues))

	return
}

// getSettingsParamters returns settings parameters of context
func getSettingsParamters(contextName string,
	config map[string]map[string]string) (settingsParameters contextTypes.SettingsParameters, err error) {

	configuration, err := contextHelpers.GetConfiguration(config)
	if err != nil {
		return contextTypes.SettingsParameters{}, err
	}

	settingsParameters = contextTypes.SettingsParameters{
		Name: contextName,
		Description: contextTypes.SettingsParametersDescription{
			Maintainer: fmt.Sprintf("\"%s\"", config["description"]["maintainer"]),
			Created:    fmt.Sprintf("\"%s\"", "2019-01-01")},
		Docker: contextTypes.SettingsParametersDocker{
			ImagesPrefix: fmt.Sprintf("%s-%s", configuration["docker"]["images-prefix"], contextName),
			Network:      fmt.Sprintf("%s-%s", configuration["docker"]["network-prefix"], contextName)},
		Paths: contextTypes.SettingsParametersPaths{
			Templates: config["paths"]["context-templates"]},
		Building: contextTypes.SettingsParametersBuilding{
			Template: configuration["building"]["template"]},
		Deploying: contextTypes.SettingsParametersDeploying{
			Template: configuration["deploying"]["template"]},
		SystemServices: contextTypes.SettingsParametersSystemServices{
			Template: configuration["system-services"]["template"]},
		ApplicationServices: contextTypes.SettingsParametersApplicationServices{
			BaseBranch:          configuration["application-services"]["base-branch"],
			FeatureBranchNaming: configuration["application-services"]["feature-branch-naming"],
			Template:            configuration["application-services"]["template"]}}

	return
}
