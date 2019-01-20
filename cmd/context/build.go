package context

import (
	contextHelpers "devlab/cmd/context/common/helpers"
	"devlab/lib/exec"
	"devlab/lib/files"
	"devlab/lib/logger"
	"errors"
	"fmt"
	"path"
	"strings"
)

// BuildService builds docker image for one application service of Context
func BuildService(contextName, serviceName string) (err error) {
	context, config, applicationServices, err := initContextToBuild(contextName)
	if err != nil {
		return
	}

	buildingSettingsPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/building.settings.yml"
	buildingSettings, err := files.ReadTwoLevelYaml(buildingSettingsPath)
	if err != nil {
		return
	}

	if applicationServices[serviceName]["enabled"] != "true" {
		return
	}

	imagesPrefix := context["docker"]["images-prefix"]
	serviceImageName := fmt.Sprintf("%s/%s", imagesPrefix, serviceName)
	isServiceImageExists, err := checkImageExists(serviceImageName)

	if isServiceImageExists {
		logger.Info("%s/%s: docker image exists", contextName, serviceName)
		return
	}

	buildingServiceParams := buildingSettings[serviceName]
	err = createBuildingServiceFolder(contextName, serviceName, context, config, buildingServiceParams)
	if err != nil {
		return
	}

	err = createBuildingImageIfNeed(contextName, serviceName, buildingServiceParams)
	if err != nil {
		return
	}

	err = runBuildingImageAndInstallDependencies(contextName, serviceName,
		config, applicationServices, buildingServiceParams)
	if err != nil {
		return
	}

	// err = buildServiceImage(contextName, serviceName, buildingServiceParams)
	// if err != nil { return }

	return
}

/************************************* helpers ************************************/

///////////////////
//initContextToBuild initializes and returns main context parameters for Context building
func initContextToBuild(contextName string) (context,
	config, applicationServices map[string]map[string]string, err error) {

	config, err = contextHelpers.GetMainConfig()
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

	applicationServices, err = contextHelpers.GetApplicationServices(contextName, context, config, configuration)
	if err != nil {
		return
	}

	context, err = files.ReadTwoLevelYaml(contextSettingsPath)

	return
}

///////////////////
// createBuildingServiceFolder checks if building service assets folder created
// and creates it if not
func createBuildingServiceFolder(contextName, serviceName string,
	context, config map[string]map[string]string,
	buildingServiceParams map[string]string) (err error) {

	//check and refresh service building folder
	buildingServiceFolderPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/building/" + serviceName
	isBuildingServiceFolderExists, _ := files.IsExists(buildingServiceFolderPath)
	if isBuildingServiceFolderExists {
		logger.Warn("%s/%s: \n     There is building folder for service: %s. \n     This folder will be recreated!",
			contextName, serviceName, buildingServiceFolderPath)
		files.Delete(buildingServiceFolderPath)
	}

	err = files.CreateDir(buildingServiceFolderPath)
	if err != nil {
		return
	}

	// copy building install script
	if err = checkBuildingInstallScript(contextName, serviceName, buildingServiceParams); err != nil {
		return
	}

	buildingInstallScriptFilename := path.Base(buildingServiceParams["install-script"])
	srcInstallScriptPath := buildingServiceParams["building-scripts-path"] + "/" + buildingServiceParams["install-script"]
	distInstallScriptPath := buildingServiceFolderPath + "/" + buildingInstallScriptFilename
	err = files.Copy(srcInstallScriptPath, distInstallScriptPath)
	if err != nil {
		return
	}

	// return if mount-ssh and set-building-envs params are set in false
	if buildingServiceParams["mount-ssh"] != "true" &&
		buildingServiceParams["set-building-envs"] != "true" {
		return
	}

	// copy ssh credential files
	if buildingServiceParams["mount-ssh"] == "true" {
		if err = checkSSHCredentials(context); err != nil {
			return
		}

		err = files.CreateDir(buildingServiceFolderPath + "/ssh")
		if err != nil {
			return
		}

		privateSSHKeyFilename := path.Base(context["ssh"]["private-key-path"])
		distPrivateSSHKeyPath := buildingServiceFolderPath + "/ssh/" + privateSSHKeyFilename
		err = files.Copy(context["ssh"]["private-key-path"], distPrivateSSHKeyPath)
		if err != nil {
			return
		}

		sshConfigFilename := path.Base(context["ssh"]["ssh-config"])
		distSSHConfigPath := buildingServiceFolderPath + "/ssh/" + sshConfigFilename
		err = files.Copy(context["ssh"]["ssh-config"], distSSHConfigPath)
		if err != nil {
			return
		}
	}

	// copy building env file
	if buildingServiceParams["set-building-envs"] == "true" {
		if err = checkBuildingEnvFile(contextName, serviceName, buildingServiceParams); err != nil {
			return
		}

		buildingEnvFileFilename := path.Base(buildingServiceParams["building-env-file"])
		srcBuildingEnvFilePath := buildingServiceParams["building-scripts-path"] + "/" + buildingServiceParams["building-env-file"]
		distBuildingEnvFilePath := buildingServiceFolderPath + "/" + buildingEnvFileFilename
		err = files.Copy(srcBuildingEnvFilePath, distBuildingEnvFilePath)
		if err != nil {
			return
		}
	}

	return
}

///////////////////
// checkSshCredentials checks ssh private key and ssh.config
func checkSSHCredentials(context map[string]map[string]string) (err error) {
	if context["ssh"]["private-key-path"] == "" {
		return errors.New("ssh private-key is not defined in context.settings.yml")
	}

	if context["ssh"]["ssh-config"] == "" {
		return errors.New("ssh config is not defined in context.settings.yml")
	}

	isSSHPrivateKeyExists, _ := files.IsExists(context["ssh"]["private-key-path"])
	if !isSSHPrivateKeyExists {
		return errors.New("ssh private-key does not exist")
	}

	isSSHConfigExists, _ := files.IsExists(context["ssh"]["ssh-config"])
	if !isSSHConfigExists {
		return errors.New("ssh config does not exist")
	}

	return
}

///////////////////
// checkBuildingEnvFile checks building env file
func checkBuildingEnvFile(contextName, serviceName string,
	buildingServiceParams map[string]string) (err error) {

	if buildingServiceParams["building-env-file"] == "" ||
		buildingServiceParams["building-scripts-path"] == "" {
		errInfo := fmt.Sprintf("%s/%s: building env file is not defined in building.settings.yml", contextName, serviceName)
		return errors.New(errInfo)
	}

	checkPath := buildingServiceParams["building-scripts-path"] + "/" + buildingServiceParams["building-env-file"]
	isBuildingEnvFileExists, _ := files.IsExists(checkPath)
	if !isBuildingEnvFileExists {
		errInfo := fmt.Sprintf("%s/%s: building env file does not exist", contextName, serviceName)
		return errors.New(errInfo)
	}

	return
}

///////////////////
// checkBuildingEnvFile checks building env file
func checkBuildingInstallScript(contextName, serviceName string,
	buildingServiceParams map[string]string) (err error) {

	if buildingServiceParams["install-script"] == "" ||
		buildingServiceParams["building-scripts-path"] == "" {
		errInfo := fmt.Sprintf("%s/%s: building install script is not defined in building.settings.yml", contextName, serviceName)
		return errors.New(errInfo)
	}

	checkPath := buildingServiceParams["building-scripts-path"] + "/" + buildingServiceParams["install-script"]
	isBuildingInstallScriptExists, _ := files.IsExists(checkPath)
	if !isBuildingInstallScriptExists {
		errInfo := fmt.Sprintf("%s/%s: building install script does not exist", contextName, serviceName)
		return errors.New(errInfo)
	}
	return
}

///////////////////
// createBuildingImageIfNeed checks if building image exists and createsit if need
func createBuildingImageIfNeed(contextName, serviceName string,
	buildingServiceParams map[string]string) (err error) {

	buildingImageName := buildingServiceParams["building-image"]
	if buildingImageName == "" {
		errInfo := fmt.Sprintf("%s/%s: building image name is not defined", contextName, serviceName)
		return errors.New(errInfo)
	}

	isBuildingImageExists, err := checkImageExists(buildingImageName)
	if err != nil {
		return
	}

	if !isBuildingImageExists {
		buildingDockerfile := buildingServiceParams["building-dockerfiles-path"] + "/" + buildingServiceParams["building-dockerfile"]

		isBuildingDockerfileExists, _ := files.IsExists(buildingDockerfile)
		if !isBuildingDockerfileExists {
			errInfo := fmt.Sprintf("%s/%s: building dockerfile does not exist: %s",
				contextName, serviceName, buildingDockerfile)
			return errors.New(errInfo)
		}

		logger.Info("Building docker image '%s'", buildingImageName)
		dockerBuildImageCmd := fmt.Sprintf("docker build -f ./%s -t %s .", buildingDockerfile, buildingImageName)
		err = exec.CommandToStdout(dockerBuildImageCmd)
	}

	return
}

///////////////////
// checkImageExists checks if docker image with name imageName exists
func checkImageExists(imageName string) (result bool, err error) {
	dockerCmd := fmt.Sprintf("docker images | grep %s | awk '{ print $3 }'", imageName)
	imagesIDs, err := exec.Command(dockerCmd)
	result = strings.TrimSpace(imagesIDs) != ""
	return
}

/////////////////
// runBuildingImageAndInstallDependencies runs docker container using building docker image
// and installs application dependencies in them.
// Installed applications dependencies volumes mounts on local machine
func runBuildingImageAndInstallDependencies(contextName, serviceName string,
	config, applicationServices map[string]map[string]string,
	buildingServiceParams map[string]string) (err error) {

	// serviceBranch := applicationServices[serviceName]["branch"]
	// buildingServiceFolderPath := "./" + config["paths"]["contexts"] + "/" + contextName + "/building/" + serviceName
	// buildingSSHFolderPath, _ := files.AbsolutePath(buildingServiceFolderPath + "/ssh")
	// servicePath, _ := files.AbsolutePath("./" + config["paths"]["contexts"] + "/" + contextName + "/services/" + serviceName)
	// installScriptPath, _ := files.AbsolutePath(buildingServiceFolderPath + "/install")
	// buildingImageName := buildingServiceParams["building-image"]

	// dockerRunInstallServiceCmd :=
	// 	fmt.Sprintf(`docker run -it --rm -e "BRANCH=%s" -v %s/:/root/.ssh/ -v %s/:/usr/src/app/ -v %s:/usr/src/app/install %s`,
	// 		serviceBranch, buildingSSHFolderPath, servicePath, installScriptPath, buildingImageName)

	//logger.Debug("dockerRunInstallServiceCmd: ", dockerRunInstallServiceCmd)
	// err = exec.CommandToStdout(dockerRunInstallServiceCmd)
	err = exec.CommandToStdoutDebug()

	return
}

// // buildServiceImage builds docker service image
// func buildServiceImage(contextName, serviceName string,
// 	buildingServiceParams map[string]map[string]string) (err error) {
// 	...
//  }
