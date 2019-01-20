package context

type ContextSettingsDescription struct {
	Maintainer string
	Created    string
}

type ContextSettingsDocker struct {
	ImagesPrefix         string
	Network              string
	RegistryHost         string
	ImagesPushPrefix     string
	DockerComposeVersion string
}

type ContextSettingsGit struct {
	ServerHost string
	AutoPush   string
}

type ContextSettingsSsh struct {
	PrivateKeyPath string
	SshConfig      string
}

type ContextSettingsPaths struct {
	Templates string
}

type ContextSettingsBuilding struct {
	Configure               string
	BuildingScriptsPath     string
	BuildingDockerfilesPath string
	Template                string
}

type ContextSettingsDeploying struct {
	Configure           string
	SystemServices      string
	ApplicationServices string
	Template            string
}

type ContextSettingsSystemServices struct {
	Configure string
	Template  string
}

type ContextSettingsApplicationServices struct {
	Configure                        string
	BaseBranch                       string
	FeatureBranchNaming              string
	MountSourceCodeVolumeOnDeploying string
	DockerRegistryTag                string
	Template                         string
}

type ContextSettings struct {
	Name                string
	Description         ContextSettingsDescription
	Git                 ContextSettingsGit
	Ssh                 ContextSettingsSsh
	Docker              ContextSettingsDocker
	Paths               ContextSettingsPaths
	Building            ContextSettingsBuilding
	Deploying           ContextSettingsDeploying
	SystemServices      ContextSettingsSystemServices
	ApplicationServices ContextSettingsApplicationServices
}
