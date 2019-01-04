package Context

import (
  "errors"
  "devlab/lib/files"
  "devlab/lib/services"
  "devlab/lib/logger"
  "strings"
)

/**
* Installing context:
* - pull all repositories of application services
* - create docker-compose.yml files or helm charts for deploying
* - building or pulling images (if it's need)
*/
func Install(contextName string) (err error) {
	if contextName == "" { return ErrNotDefinedContextName }

  config, err := files.ReadMainConfig()
  if err != nil { return ErrCouldntReadConfig }   
  if config["paths"]["contexts"] == "" { return ErrNotDefinedContextsPath } 

  // Check context dir and settings.yml and create it if need  
	contextDir := "./" + config["paths"]["contexts"] + "/" + contextName 
	contextSettingsPath := contextDir + "/context.settings.yml"
	isContextDirExists, _ :=  files.IsExists(contextDir)
	isContextSettingsExists, _ := files.IsExists(contextSettingsPath)
  if !isContextDirExists || !isContextSettingsExists {		
    logger.Warn(`Context '%s' is not created! 
   New context will be created: %s. 
   The installing process has been stoped.  
   You should configure new context and start installation again.`, contextName, contextSettingsPath)

		if err = Create(contextName); err != nil { return err }

		return ErrContextIsNotCreated
  }
  
  configuration, err := getConfiguration(config)
  if err != nil { return err }

  // Read context settings
  context, err := files.ReadTwoLevelYaml(contextSettingsPath)
  if err != nil { return err }

  // Check context services dir and create it if need
  contextServicesDir := config["paths"]["contexts"] + "/" + contextName + "/services"
  isContextServicesDirExists, err :=  files.IsExists("./" + contextServicesDir)
  if err != nil { return err }
  if !isContextServicesDirExists {
    files.CreateDir("./" + contextServicesDir)
  }

  // Set task base branch
  taskBaseBranch := getValueFromContextOrDefault(context, configuration, "git", "base-branch")

  // Set template path
  templatesPath := getValueFromContextOrDefault(context, config, "paths", "context-templates")

  // Set application services template
  applicationServicesTemplate := getValueFromContextOrDefault(context, configuration, "application-services", "template") 

  // Check services repo and clone/refresh it if need
  applicationServicesTemplatePath := "./" + templatesPath + "/application-services/" + applicationServicesTemplate
  if isAplicationServicesDirExists, _ :=  files.IsExists(applicationServicesTemplatePath); !isAplicationServicesDirExists {
    return errors.New("ERROR: applications services directory does not exist!")
  } 

  applicationServices, err := files.ReadTwoLevelYaml(applicationServicesTemplatePath)
  if err != nil { return err }

  for serviceName, serviceParams := range applicationServices { 
    logger.Header(strings.ToUpper(serviceName))

    isServiceDirExists, _ :=  files.IsExists("./" + contextServicesDir + "/serviceName") 
        
    if !isServiceDirExists {      
      services.Clone(contextServicesDir, serviceName, config["git"]["git-server-host"], serviceParams["github-path"] )
    }

    serviceBaseBranch := serviceParams["base-branch"]    
    if  serviceBaseBranch == "" {
      serviceBaseBranch = taskBaseBranch
    }

    contextServiceBranch := serviceParams["branch"]
    if  contextServiceBranch == "" {
      contextServiceBranch = serviceBaseBranch
    }

    services.RefreshGitRepo(contextServicesDir, serviceName, contextServiceBranch, serviceBaseBranch)
  }

  return
}