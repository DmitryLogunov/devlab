package context

import (
	"devlab/lib/files"
	"devlab/lib/logger"
	"devlab/lib/services"
	"strings"

	contextErrors "devlab/cmd/context/common/errors"
	contextHelpers "devlab/cmd/context/common/helpers"
)

// Install context:
//  - pull all git repositories of application services and refresh them branches
//  - building or pulling images (if it's need)
//  - create docker-compose.yml files or helm charts for deploying
func Install(contextName string) (err error) {
	config, configuration, context, contextServicesDir, err := initContextToInstall(contextName)
	if err != nil {
		return
	}

	applicationServices, err := contextHelpers.GetApplicationServices(contextName, context, config, configuration)
	if err != nil {
		return
	}

	// clone and refresh all application-services repositories from git server
	taskBaseBranch := contextHelpers.GetValueFromContextOrDefault(context, configuration,
		"application-services", "base-branch")
	err = cloneAndRefreshApplicationServicesGitRepo(contextServicesDir, taskBaseBranch,
		context, applicationServices)

	for serviceName := range applicationServices {
		if err = BuildService(contextName, serviceName); err != nil {
			break
		}
	}

	return
}

// @private

// initContextToCreat initializes and returns main context parameters for  Context installation
func initContextToInstall(contextName string) (config,
	configuration, context map[string]map[string]string,
	contextServicesDir string, err error) {

	if contextName == "" {
		err = contextErrors.ErrNotDefinedContextName
		return
	}

	config, err = contextHelpers.GetMainConfig()
	if err != nil {
		return
	}

	contextSettingsPath, err := checkAndCreateContextSettingsIfNotExists(contextName, config)
	if err != nil {
		return
	}

	configuration, err = contextHelpers.GetConfiguration(config)
	if err != nil {
		return
	}

	context, err = files.ReadTwoLevelYaml(contextSettingsPath)
	if err != nil {
		return
	}

	contextServicesDir, err = getContextServicesDir(contextName, config)

	return
}

// checkAndCreateContextSettingsIfNotExists checks context dir and settings.yml and creates it if need
func checkAndCreateContextSettingsIfNotExists(contextName string,
	config map[string]map[string]string) (contextSettingsPath string, err error) {

	contextDir := "./" + config["paths"]["contexts"] + "/" + contextName
	contextSettingsPath = contextDir + "/context.settings.yml"
	isContextDirExists, _ := files.IsExists(contextDir)
	isContextSettingsExists, _ := files.IsExists(contextSettingsPath)
	if !isContextDirExists || !isContextSettingsExists {
		logger.Warn(`Context '%s' is not created! 
   New context will be created: %s. 
   The installing process has been stoped.  
   You should configure new context and start installation again.\n`, contextName, contextDir)

		if err = Create(contextName); err != nil {
			return
		}

		return contextSettingsPath, contextErrors.ErrContextIsNotCreated
	}
	return
}

// getContextServicesDir returns context servcies directory
func getContextServicesDir(contextName string, config map[string]map[string]string) (contextServicesDir string, err error) {
	contextServicesDir = config["paths"]["contexts"] + "/" + contextName + "/services"
	isContextServicesDirExists, err := files.IsExists("./" + contextServicesDir)
	if err != nil {
		return
	}
	if !isContextServicesDirExists {
		files.CreateDir("./" + contextServicesDir)
	}

	return
}

// cloneAndRefreshApplicationServicesGitRepo clones and refreshs all application-services repositories
// from git server
func cloneAndRefreshApplicationServicesGitRepo(contextServicesDir, taskBaseBranch string,
	context, applicationServices map[string]map[string]string) (err error) {

	for serviceName, serviceParams := range applicationServices {
		if serviceParams["enabled"] != "true" {
			continue
		}

		logger.Header(strings.ToUpper(serviceName))

		isServiceDirExists, _ := files.IsExists("./" + contextServicesDir + "/serviceName")

		if !isServiceDirExists {
			services.Clone(contextServicesDir,
				serviceName,
				context["git"]["server-host"],
				serviceParams["git-path"])
		}

		baseServiceBranch := serviceParams["base-branch"]
		if baseServiceBranch == "" {
			baseServiceBranch = taskBaseBranch
		}

		contextServiceBranch := serviceParams["branch"]
		if contextServiceBranch == "" {
			contextServiceBranch = baseServiceBranch
		}

		services.RefreshGitRepo(contextServicesDir, serviceName,
			contextServiceBranch, baseServiceBranch, context)
	}

	return
}
