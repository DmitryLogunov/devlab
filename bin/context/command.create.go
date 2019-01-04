package Context

import (
  "fmt"
  "errors"
  "devlab/lib/files"
  "devlab/lib/logger"
)

/**
*  Creating a new context
*  @param contextName {string} - context name
*/
func Create(contextName string) (err error) {
  if contextName == "" { return ErrNotDefinedContextName }

  config, err := files.ReadMainConfig()
  if err != nil { return ErrCouldntReadConfig }   
  if config["paths"]["contexts"] == "" { return ErrNotDefinedContextsPath } 

  // Check context dir and create it if need  
  contextDir := "./" + config["paths"]["contexts"] + "/" + contextName 
  if isContextDirExists, _ :=  files.IsExists(contextDir); !isContextDirExists {
    files.CreateDir("./" + contextDir)
  }
  
  settingsParameters, err := getSettingsParamters(contextName, config)
  if err != nil { return err }  
 
  if config["paths"]["context-templates"] == "" { 
    return ErrNotDefinedContextsPath 
  } 
  defaultSettingsPath := "./" + config["paths"]["context-templates"] + "/default-settings.yml"

  contextSettingsPath := contextDir + "/context.settings.yml"
  isContextSettingsExists, _ :=  files.IsExists(contextSettingsPath)
  if isContextSettingsExists {
     logger.Warn("There is settings.yml file in context: %s. This file will be rewrited!", contextSettingsPath)
     files.Delete(contextSettingsPath) 
  }

  err = files.RenderTextTemplate(defaultSettingsPath, contextSettingsPath, settingsParameters)  
  if err != nil { return err }  

  err = createContextBlockSettingsFile(contextName, config, "building")
  if err != nil { return err }

  err = createContextBlockSettingsFile(contextName, config, "deploying")
  if err != nil { return err }

  err = createContextBlockSettingsFile(contextName, config, "system-services")
  if err != nil { return err }

  err = createContextBlockSettingsFile(contextName, config, "application-services")
  if err != nil { return err }

  return
}


/**
* Getting settings parameters of context
*/
func getSettingsParamters(contextName string, config map[string]map[string]string) (settingsParameters SettingsParameters, err error) {
  configuration, err := getConfiguration(config)
  if err != nil {
    return SettingsParameters{}, err
  }

  settingsParameters = SettingsParameters{
    Name: contextName,
    Description: SettingsParametersDescription{
      Maintainer: fmt.Sprintf("\"%s\"", config["description"]["maintainer"]),
      Created: fmt.Sprintf("\"%s\"", "2019-01-01")},
    Git: SettingsParametersGit{
      BaseBranch: configuration["git"]["base-branch"]},
    Docker: SettingsParametersDocker{
      ImagesPrefix: fmt.Sprintf("%s-%s", configuration["docker"]["images-prefix"], contextName),
      Network: fmt.Sprintf("%s-%s", configuration["docker"]["network-prefix"], contextName)},
    Paths: SettingsParametersPaths{
      Templates: config["paths"]["context-templates"]},
    Building: SettingsParametersBuilding{
      Template: configuration["building"]["template"]},
    Deploying: SettingsParametersDeploying{
      Template: configuration["deploying"]["template"]},
    SystemServices: SettingsParametersSystemServices{
      Template: configuration["system-services"]["template"]},
    ApplicationServices: SettingsParametersApplicationServices{
      Template: configuration["application-services"]["template"]}}

  return
}

/**
*  Create block of context settings yaml file 
*/
func createContextBlockSettingsFile(contextName string, config map[string]map[string]string, blockName string) (err error) {
  configuration, err := getConfiguration(config)
  if err != nil { return err }

  if configuration[blockName]["template"] == "" {
    errInfo := fmt.Sprintf("ERROR: context => Template for block '%s' is not defined. Couldn't create file settings for this block", blockName)
    return errors.New(errInfo)    
  }

  // Read context settings
	contextSettingsPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/context.settings.yml"
  context, err := files.ReadTwoLevelYaml(contextSettingsPath)
  if err != nil { return err }

  // Set template path
  templatesPath := getValueFromContextOrDefault(context, config, "paths", "context-templates")

  // Set application services template
  contextBlockTemplate := getValueFromContextOrDefault(context, configuration, blockName, "template")
  contextBlockTemplatePath := "./" + templatesPath + "/" + blockName + "/" + contextBlockTemplate
  contextBlock, err := files.ReadTwoLevelYaml(contextBlockTemplatePath)
  if err != nil { return err }
  
  // Read parameters of context block
  parametersPath := "./" + templatesPath + "/" + blockName + "/parameters.yml"
  parametersOfBlockItems, err := files.ReadOneLevelYaml(parametersPath)
  if err != nil { return err }   
  
  // Check if context block settings file exists and delete it if yes
  contextBlockSettingsFilePath := "./" + config["paths"]["contexts"] + "/" + contextName + "/" + blockName + ".settings.yml"
  isCcontextBlockSettingsFileExists, _ :=  files.IsExists(contextBlockSettingsFilePath)
  if isCcontextBlockSettingsFileExists {
     files.Delete(contextBlockSettingsFilePath) 
  }

  // Fill context block settings map
  contextBlockSettings := make(map[string]map[string]string)
  for item, _ := range contextBlock { 
    contextBlockSettings[item] = parametersOfBlockItems
  }

  // Create context block settings yaml file
  err = files.WriteYaml(contextBlockSettingsFilePath, contextBlockSettings)
  if err != nil { return err } 

  return
}
