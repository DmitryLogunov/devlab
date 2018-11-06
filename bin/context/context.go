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

  context, err := files.ReadContextConfig("./" + config["contexts-path"] + "/" + contextName + "/settings.yml")
  if errors.CheckAndReturnIfError(err) { return }

  contextServicesDir := config["contexts-path"] + "/" + contextName + "/services"

  isContextServicesDirExists, err :=  files.IsExists("./" + contextServicesDir)
  if errors.CheckAndReturnIfError(err) { return }

  if !isContextServicesDirExists {
    files.CreateDir("./" + contextServicesDir)
  }

  taskBaseBranch := context["context"]["task"]["base-branch"]
  if taskBaseBranch == "" {
    taskBaseBranch = config["base-branch"]
  }

  for serviceName, serviceParams := range context["application-services"] { 
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
  

