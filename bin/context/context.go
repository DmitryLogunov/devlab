package Context

import (
  "fmt"
  "devlab/lib/files"
  "devlab/lib/errors"
)

/**
*
*/
func Set(contextName string) (err error) {
  config, err := files.ReadMainConfig()
  if errors.CheckAndReturnIfError(err) { return }

  // Check context dir and create it if need  
  contextDir := "./" + config["paths"]["contexts"] + "/" + contextName 
  isContextDirExists, _ :=  files.IsExists("./" + contextDir)
  if !isContextDirExists {
    files.CreateDir("./" + contextDir)
  }
  
  settingsParameters := getSettingsParamters(contextName, config)  
 
  defaultSettingsPath := "./" + config["paths"]["contexts-templates"] + "/default-settings.yml"

  contextSettingsPath := contextDir + "/settings.yml"
  isContextSettingsExists, _ :=  files.IsExists(contextSettingsPath)
  if isContextSettingsExists {
     files.Delete(contextSettingsPath) 
  }
  
  files.RenderTextTemplate(defaultSettingsPath, contextSettingsPath, settingsParameters)

  return 
}

/**
*
*/
func  getSettingsParamters(contextName string, config map[string]map[string]string) SettingsParameters {
  configurationTemplatesPath := config["paths"]["configurations"] + "/" + config["description"]["configuration"]
  isConfigurationExists, _ :=  files.IsExists("./" + configurationTemplatesPath)
  if !isConfigurationExists { return SettingsParameters{} }
  configuration, _ := files.ReadTwoLevelYaml(configurationTemplatesPath)

  settingsParameters := SettingsParameters{
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
      Templates: config["paths"]["contexts-templates"]},
    Building: SettingsParametersBuilding{
      Template: configuration["building"]["template"]},
    Deploying: SettingsParametersDeploying{
      Template: configuration["deploying"]["template"]},
    SystemServices: SettingsParametersSystemServices{
      Template: configuration["system-services"]["template"]},
    ApplicatonServices: SettingsParametersApplicatonServices{
      Template: configuration["applicaton-services"]["template"]}}

  return settingsParameters
}