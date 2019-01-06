package context

import (
	"devlab/lib/files"
	"devlab/lib/logger"
	"devlab/lib/services"
	"errors"
	"strings"

	contextErrors "devlab/cmd/context/common/errors"
	contextHelpers "devlab/cmd/context/common/helpers"
)

// Install context:
//  - pull all git repositories of application services and refresh them branches
//  - create docker-compose.yml files or helm charts for deploying
//  - building or pulling images (if it's need)
func Install(contextName string) (err error) {
	if contextName == "" {
		return contextErrors.ErrNotDefinedContextName
	}

	config, err := getMainConfig()
	if err != nil {
		return
	}

	contextSettingsPath, err := checkAndCreateContextSettingsIfNotExists(contextName, config)
	if err != nil {
		return
	}

	configuration, err := contextHelpers.GetConfiguration(config)
	if err != nil {
		return
	}

	context, err := files.ReadTwoLevelYaml(contextSettingsPath)
	if err != nil {
		return
	}

	contextServicesDir, err := getContextServicesDir(contextName, config)
	if err != nil {
		return
	}

	applicationServices, err := getApplicationServices(contextName, context, config, configuration)
	if err != nil {
		return
	}

	// clone and refresh all application-services from git server
	taskBaseBranch := contextHelpers.GetValueFromContextOrDefault(context, configuration, "application-services", "base-branch")
	err = cloneAndRefreshApplicationServicesGitRepo(contextServicesDir, taskBaseBranch, config, applicationServices)

	return
}

/***************************************************************/

// Get main config
func getMainConfig() (config map[string]map[string]string, err error) {
	config, err = files.ReadMainConfig()
	if err != nil {
		return make(map[string]map[string]string), contextErrors.ErrCouldntReadConfig
	}
	if config["paths"]["contexts"] == "" {
		return make(map[string]map[string]string), contextErrors.ErrNotDefinedContextsPath
	}
	return
}

// Check context dir and settings.yml and create it if need
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
   You should configure new context and start installation again.`, contextName, contextSettingsPath)

		if err = Create(contextName); err != nil {
			return
		}

		return contextSettingsPath, contextErrors.ErrContextIsNotCreated
	}
	return
}

// Getting application services settings by merging default template settings and context level settings
func getApplicationServices(contextName string, context, config map[string]map[string]string,
	configuration map[string]map[string]string) (applicationServices map[string]map[string]string, err error) {

	applicationServices = make(map[string]map[string]string)

	templatesPath := contextHelpers.GetValueFromContextOrDefault(context, config, "paths", "context-templates")
	applicationServicesTemplate := contextHelpers.GetValueFromContextOrDefault(context, configuration, "application-services", "template")
	applicationServicesTemplatePath := "./" + templatesPath + "/application-services/" + applicationServicesTemplate

	isAplicationServicesDirExists, _ := files.IsExists(applicationServicesTemplatePath)
	if !isAplicationServicesDirExists {
		err = errors.New("applications services directory does not exist")
		return
	}

	applicationServicesFromTemplate, err := files.ReadTwoLevelYaml(applicationServicesTemplatePath)
	if err != nil {
		return
	}

	// applications services settings equal template settings as default
	applicationServices = applicationServicesFromTemplate

	// checking if application-services settings from context exist and merge its if yes
	applicationServicesContextSettingsPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/application-services.settings.yml"
	isAplicationServicesContextSettingsExists, _ := files.IsExists(applicationServicesContextSettingsPath)
	if isAplicationServicesContextSettingsExists {
		applicationServicesFromContext, _ := files.ReadTwoLevelYaml(applicationServicesContextSettingsPath)
		applicationServices = contextHelpers.MergeMaps(applicationServicesFromContext, applicationServicesFromTemplate)
	}

	return
}

// Get context servcies directory
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

// Clone and refresh all application-services from git server
func cloneAndRefreshApplicationServicesGitRepo(contextServicesDir, taskBaseBranch string,
	config, applicationServices map[string]map[string]string) (err error) {

	for serviceName, serviceParams := range applicationServices {
		logger.Header(strings.ToUpper(serviceName))

		isServiceDirExists, _ := files.IsExists("./" + contextServicesDir + "/serviceName")

		if !isServiceDirExists {
			services.Clone(contextServicesDir, serviceName, config["git"]["git-server-host"], serviceParams["git-path"])
		}

		serviceBaseBranch := serviceParams["base-branch"]
		if serviceBaseBranch == "" {
			serviceBaseBranch = taskBaseBranch
		}

		contextServiceBranch := serviceParams["branch"]
		if contextServiceBranch == "" {
			contextServiceBranch = serviceBaseBranch
		}

		services.RefreshGitRepo(contextServicesDir, serviceName, contextServiceBranch, serviceBaseBranch)
	}

	return
}
