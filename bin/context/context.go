package Context

import (
  "devlab/lib/logger"
  "devlab/lib/files"
  "devlab/lib/errors"
  "devlab/lib/services"
  "strings"
)


func Set(contextName string) (err error) {
  config, err := files.ReadMainConfig()
  if errors.CheckAndReturnIfError(err) { return }

  // Check context dir and create it if need  
  contextDir := "./" + config["contexts-path"] + "/" + contextName 
  isContextDirExists, _ :=  files.IsExists("./" + contextDir)
  if !isContextDirExists {
    files.CreateDir("./" + contextDir)
  }
  
  // Check context settings file and create it if need  
  contextSettings := contextDir + "/settings.yml"
  isContextSettingsExists, _ :=  files.IsExists(contextSettings)
  if !isContextSettingsExists {
    files.Copy("./" + config["data-path"] + "/default-context.yml", contextSettings) 
  }

  // Read context settings
  context, err := files.ReadContextConfig(contextSettings)
  if errors.CheckAndReturnIfError(err) { return }

  // Check context services dir and create it if need
  contextServicesDir := config["contexts-path"] + "/" + contextName + "/services"
  isContextServicesDirExists, err :=  files.IsExists("./" + contextServicesDir)
  if errors.CheckAndReturnIfError(err) { return }
  if !isContextServicesDirExists {
    files.CreateDir("./" + contextServicesDir)
  }

  // Set task base branch
  taskBaseBranch := context["context"]["task"]["base-branch"]
  if taskBaseBranch == "" {
    taskBaseBranch = config["base-branch"]
  }

  // Check services repo and clone/refresh it if need
  applicationServices := context["applicaton-services"]
  for serviceName, serviceParams := range applicationServices { 

    logger.Header(strings.ToUpper(serviceName))

    isServiceDirExists, _ :=  files.IsExists("./" + contextServicesDir + "/serviceName") 
        
    if !isServiceDirExists {      
      services.Clone(contextServicesDir, serviceName, config["github-repository-path"], serviceParams["github-path"] )
    }

    serviceBaseBranch := serviceParams["base-branch"]    
    if  serviceBaseBranch == "" {
      serviceBaseBranch = taskBaseBranch
    }
    services.RefreshGitRepo(contextServicesDir, serviceName, serviceParams["branch"], serviceBaseBranch)
  }
  
  return
}

func Create(contextName string) (err error) {
  return
}
  

